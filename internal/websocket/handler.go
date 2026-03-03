package websocket

import (
	"chatvis-chat/internal/domain"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

// WebSocketController gestionará las conexiones y usará el Hub
type WebSocketController struct {
	Hub          *Hub
	GrupoUseCase domain.GrupoUseCase
}

// NewWebSocketController crea un nuevo controlador de WebSocket
func NewWebSocketController(h *Hub, gu domain.GrupoUseCase) *WebSocketController {
	return &WebSocketController{Hub: h, GrupoUseCase: gu}
}

// WebSocketUpgrade es el handler que actualiza la conexión HTTP a una WebSocket
func (c *WebSocketController) WebSocketUpgrade(ctx *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(ctx) {
		return fiber.ErrUpgradeRequired
	}
	return ctx.Next()
}

// WebSocketChat handles WebSocket connections, authenticates, and subscribes users to groups.
func (c *WebSocketController) WebSocketChat(conn *websocket.Conn) {
	var authMsg struct {
		Token string `json:"token"`
	}

	if err := conn.ReadJSON(&authMsg); err != nil {
		log.Println("Error al leer el token de autenticación:", err)
		conn.Close()
		return
	}

	token, err := jwt.Parse(authMsg.Token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY_JWT")), nil
	})

	if err != nil || !token.Valid {
		log.Println("Token JWT inválido:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("authentication_failed"))
		conn.Close()
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	userIDInterface, ok := claims["id"]
	if !ok {
		log.Println("ID de usuario no válido en el token.")
		conn.WriteMessage(websocket.TextMessage, []byte("authentication_failed"))
		conn.Close()
		return
	}
	userIDStr := fmt.Sprintf("%v", userIDInterface)
	if len(strings.TrimSpace(userIDStr)) == 0 {
		log.Println("ID de usuario vacío en el token.")
		conn.WriteMessage(websocket.TextMessage, []byte("authentication_failed"))
		conn.Close()
		return
	}

	idUser, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.Println("Error al convertir el ID de usuario:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("authentication_failed"))
		conn.Close()
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("authentication_successful"))

	groupClaves, err := c.GrupoUseCase.GetAllGruposByUsuarioIdToClaves(idUser)
	if err != nil {
		log.Printf("Error al obtener grupos para el usuario %s: %v", userIDStr, err)
		conn.Close()
		return
	}

	c.Hub.Register(userIDStr, conn)
	c.Hub.SubscribeUserToGroups(userIDStr, groupClaves)

	defer func() {
		c.Hub.Unregister(userIDStr)
		conn.Close()
	}()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Error al leer mensaje de %s: %v\n", userIDStr, err)
			break
		}

		if !c.Hub.CheckUserInGroup(userIDStr, msg.GroupID) {
			log.Printf("Usuario %s no pertenece al grupo %s. Mensaje no enviado.", userIDStr, msg.GroupID)
			continue
		}

		msg.SenderID = userIDStr
		c.Hub.Broadcast(msg)
	}
}
