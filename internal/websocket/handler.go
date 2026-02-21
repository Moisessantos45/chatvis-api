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
	// 1. Read the authentication token from the first message
	var authMsg struct {
		Token string `json:"token"`
	}

	if err := conn.ReadJSON(&authMsg); err != nil {
		log.Println("Error al leer el token de autenticación:", err)
		conn.Close()
		return
	}

	// 2. Validate the token and get the user ID
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

	// Get the user ID from claims. The claim key is "id" (lowercase as decoded by standard).
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

	// Convert the string ID to uint64 for the service call
	idUser, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.Println("Error al convertir el ID de usuario:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("authentication_failed"))
		conn.Close()
		return
	}

	// 3. Confirm authentication success and continue
	conn.WriteMessage(websocket.TextMessage, []byte("authentication_successful"))

	// 4. Subscribe the user to their groups
	groupClaves, err := c.GrupoUseCase.GetAllGruposByUsuarioIdToClaves(idUser)
	if err != nil {
		log.Printf("Error al obtener grupos para el usuario %s: %v", userIDStr, err)
		conn.Close()
		return
	}

	// Register user with the string ID
	c.Hub.Register(userIDStr, conn)
	c.Hub.SubscribeUserToGroups(userIDStr, groupClaves)

	defer func() {
		c.Hub.Unregister(userIDStr)
		conn.Close()
	}()

	// 5. Handle subsequent chat messages
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
