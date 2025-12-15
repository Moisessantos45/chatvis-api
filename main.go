package main

import (
	"chatvis-chat/internal/app/auth/handleauth"
	"chatvis-chat/internal/app/auth/serviceauth"
	"chatvis-chat/internal/app/grupo"
	"chatvis-chat/internal/app/grupo/handlegrupo"
	"chatvis-chat/internal/app/grupo/servicegrupo"
	"chatvis-chat/internal/app/grupousario"
	"chatvis-chat/internal/app/grupousario/handlegrupousuario"
	"chatvis-chat/internal/app/grupousario/servicegrupousuario"
	"chatvis-chat/internal/app/ia"
	"chatvis-chat/internal/app/mensaje"
	"chatvis-chat/internal/app/mensaje/handlemensaje"
	"chatvis-chat/internal/app/mensaje/servicemensaje"
	"chatvis-chat/internal/app/routers"
	"chatvis-chat/internal/app/usuario"
	"chatvis-chat/internal/app/usuario/handleusuario"
	"chatvis-chat/internal/app/usuario/serviceusuario"
	"chatvis-chat/internal/app/websocket/handle"
	"chatvis-chat/internal/app/websocket/hub"
	"chatvis-chat/internal/db"
	"chatvis-chat/internal/pkg/middleware"
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

	grupoUsuarios := grupousario.GruposUsuariosRepository{DB: db.DB}

	grupoRepository := grupo.GrupoRepository{DB: db.DB, RepoRelaciones: &grupoUsuarios}
	grupoService := servicegrupo.NewGrupoService(grupoRepository)
	grupoController := &handlegrupo.GrupoController{Service: grupoService}

	usuarioRepository := usuario.UsuarioRepository{DB: db.DB}
	usuarioService := serviceusuario.NewUsuarioUseCase(usuarioRepository)
	usuarioController := &handleusuario.UsuarioController{Service: usuarioService}

	mensajeRepository := mensaje.MensajeRepository{DB: db.DB}
	mensajeService := servicemensaje.NewMensajeService(mensajeRepository)
	mensajeController := &handlemensaje.MensajeController{Service: mensajeService}

	autservice := serviceauth.NewServiceAuth(usuarioRepository)
	authController := &handleauth.AuthController{
		Service: autservice,
	}

	// Inicialización del Hub y el controlador de WebSocket
	wsHub := hub.NewHub()
	go wsHub.Run()

	// --- Inicialización de los servicios de IA ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	aiServices := make(map[string]*ia.AIService)
	for _, config := range ia.AiConfigurations {
		aiService := ia.NewAIService(wsHub, mensajeService, grupoService, config)
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
				if wsHub.CheckUserInGroup(service.UserID, msg.GroupID) {
					service.InputChannel() <- msg
					break
				}
			}
		}
	}()

	// Pasar el Hub Y el GrupoService al controlador de WebSocket
	webSocketController := handle.NewWebSocketController(wsHub, grupoService)

	public := app.Group("/api/public")
	routers.RegisterUsuarioRoutesBasicNotAuth(public, usuarioController)
	routers.RegisterAuthRoutes(public, authController)

	//WebSocket se define en el grupo protegido
	public.Get("/ws/chat", webSocketController.WebSocketUpgrade, websocket.New(webSocketController.WebSocketChat))

	protected := app.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())

	auth := protected.Group("/auth")
	routers.RegisterAuthRoutesWithMiddleware(auth, authController)

	grupo := protected.Group("/grupo")
	routers.RegisterGrupoRoutes(grupo, grupoController)

	usuario := protected.Group("/usuario")
	routers.RegisterUsuarioRoutes(usuario, usuarioController)

	mensaje := protected.Group("/mensaje")
	routers.RegisterMensajeRoutes(mensaje, mensajeController)

	grupoUsuario := protected.Group("/grupo-usuario")
	handleGrupoUsuario := grupousario.GruposUsuariosRepository{DB: db.DB}
	grupoUsuarioService := servicegrupousuario.NewGruposUsuariosServiceClient(handleGrupoUsuario, grupoRepository)
	grupoUsuarioController := &handlegrupousuario.GruposUsuariosController{Service: grupoUsuarioService}
	routers.RegisterGrupoUsuarioRoutes(grupoUsuario, grupoUsuarioController)

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
