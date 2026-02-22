package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"

	usuarioHttp "chatvis-chat/internal/usuario/delivery/http"
	usuarioRepo "chatvis-chat/internal/usuario/repository"
	usuarioUseCase "chatvis-chat/internal/usuario/usecase"

	mensajeHttp "chatvis-chat/internal/mensaje/delivery/http"
	mensajeRepo "chatvis-chat/internal/mensaje/repository"
	mensajeUseCase "chatvis-chat/internal/mensaje/usecase"

	grupoHttp "chatvis-chat/internal/grupo/delivery/http"
	grupoRepo "chatvis-chat/internal/grupo/repository"
	grupoUseCase "chatvis-chat/internal/grupo/usecase"

	grupoUsuarioHttp "chatvis-chat/internal/grupousuario/delivery/http"
	grupoUsuarioRepo "chatvis-chat/internal/grupousuario/repository"
	grupoUsuarioUseCase "chatvis-chat/internal/grupousuario/usecase"

	authHttp "chatvis-chat/internal/auth/delivery/http"
	authUseCase "chatvis-chat/internal/auth/usecase"

	"chatvis-chat/internal/ia"
	appWs "chatvis-chat/internal/websocket"

	"chatvis-chat/config/db"
	"chatvis-chat/internal/pkg/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
		return
	}

	err = db.Init()
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
		return
	}

	fmt.Println("¡Hola, mundo desde Go!")

	app := fiber.New()
	app.Use(compress.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "http://localhost:5173",
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Access-Control-Allow-Methods",
		MaxAge:        300,
	}))
	fmt.Println("Hello, World!")

	// ******************************
	// INYECCIÓN DE DEPENDENCIA DE USUARIO (Clean Architecture)
	pgUserRepo := usuarioRepo.NewPostgresUsuarioRepository(db.DB)
	userUseCase := usuarioUseCase.NewUsuarioUseCase(pgUserRepo)
	// ******************************

	// ******************************
	// INYECCIÓN DE DEPENDENCIA DE MENSAJE (Clean Architecture)
	pgMensajeRepo := mensajeRepo.NewPostgresMensajeRepository(db.DB)
	msgUseCase := mensajeUseCase.NewMensajeUseCase(pgMensajeRepo)
	// ******************************

	// ******************************
	// INYECCIÓN DE DEPENDENCIA DE GRUPO (Clean Architecture)
	pgGrupoRepo := grupoRepo.NewPostgresGrupoRepository(db.DB)
	grpUseCase := grupoUseCase.NewGrupoUseCase(pgGrupoRepo)
	// ******************************

	// ******************************
	// INYECCIÓN DE DEPENDENCIA DE GRUPO_USUARIO (Clean Architecture)
	pgGrupoUsuarioRepo := grupoUsuarioRepo.NewPostgresGrupoUsuarioRepository(db.DB)
	grpUsuarioUseCase := grupoUsuarioUseCase.NewGrupoUsuarioUseCase(pgGrupoUsuarioRepo, pgGrupoRepo)
	// ******************************

	// ******************************
	// INYECCIÓN DE DEPENDENCIA DE AUTH (Clean Architecture)
	authUsecase := authUseCase.NewAuthUseCase(pgUserRepo)
	// ******************************

	// Inicialización del Hub y el controlador de WebSocket
	wsHub := appWs.NewHub()
	go wsHub.Run()

	// --- Inicialización de los servicios de IA ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	enableAI := os.Getenv("ENABLE_AI_MODELS")
	aiServices := make(map[string]*ia.AIService)

	if enableAI == "true" {
		log.Println("🤖 Servicios de IA Habilitados (ENABLE_AI_MODELS=true)")
		for _, config := range ia.AiConfigurations {
			aiService := ia.NewAIService(wsHub, msgUseCase, grpUseCase, config)
			aiServices[config.UserID] = aiService
			go aiService.Start(ctx, 5) // 5 workers por servicio
		}

		// --- Enrutamiento de mensajes hacia las IA ---
		go func() {
			for msg := range wsHub.AIChannel() {
				// Ignorar mensajes de las propias IA
				if _, ok := aiServices[msg.SenderID]; ok {
					continue
				}

				// Mandar al servicio de IA correspondiente
				for _, service := range aiServices {
					if wsHub.CheckUserInGroup(service.Config.UserID, msg.GroupID) {
						service.InputChannel() <- msg
						break
					}
				}
			}
		}()
	} else {
		log.Println("🛑 Servicios de IA Deshabilitados. Cambiar ENABLE_AI_MODELS=true en .env para activar.")
		// Opcional: Podrías querer drenar el AIChannel si no se usa para evitar que se bloquee la memoria
		go func() {
			for range wsHub.AIChannel() {
				// Sink (descarte de mensajes)
			}
		}()
	}

	// Pasar el Hub Y el GrupoService al controlador de WebSocket
	webSocketController := appWs.NewWebSocketController(wsHub, grpUseCase)

	public := app.Group("/api/public")

	// Registro de rutas publicas
	usuarioHttp.NewUsuarioPublicHandler(public, userUseCase)
	authHttp.NewAuthHandler(public, authUsecase)

	//WebSocket se define en el grupo protegido
	public.Get("/ws/chat", webSocketController.WebSocketUpgrade, websocket.New(webSocketController.WebSocketChat))

	protected := app.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())

	auth := protected.Group("/auth")
	authHttp.NewAuthProtectedHandler(auth, authUsecase)

	grupo := protected.Group("/grupo")
	grupoHttp.NewGrupoHandler(grupo, grpUseCase)

	usuarioGrp := protected.Group("/usuario")
	// Registro de rutas privadas
	usuarioHttp.NewUsuarioHandler(usuarioGrp, userUseCase)

	mensajeGrp := protected.Group("/mensaje")
	mensajeHttp.NewMensajeHandler(mensajeGrp, msgUseCase)

	grupoUsuarioGrp := protected.Group("/grupo-usuario")
	grupoUsuarioHttp.NewGrupoUsuarioHandler(grupoUsuarioGrp, grpUsuarioUseCase)

	// --- Señales de cierre ---
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":3100"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	fmt.Println("Server is running on port 3100")

	<-c
	log.Println("Señal de cierre recibida. Deteniendo servicios...")

	// Cancelar el contexto -> detiene workers y listeners
	cancel()
	wsHub.Shutdown()

	// Detener Fiber con timeout
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := app.ShutdownWithContext(ctxTimeout); err != nil {
		log.Fatalf("Error al detener el servidor: %v", err)
	}

	log.Println("Aplicación detenida con éxito.")
}
