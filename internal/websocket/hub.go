package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Hub gestiona la difusión de mensajes a clientes por grupo
type Hub struct {
	clients    map[string]*websocket.Conn // ID de usuario -> Conexión
	userGroups map[string]map[string]bool // ID de usuario -> {ID de grupo: true}

	// Canales de control
	register   chan RegisterClient
	unregister chan string
	broadcast  chan Message
	aiChannel  chan Message
	done       chan struct{}

	mu sync.Mutex // Mutex para proteger el estado concurrente
}

// RegisterClient encapsula la información para el registro
type RegisterClient struct {
	UserID string
	Conn   *websocket.Conn
}

// Message representa un mensaje con el contenido y el grupo de destino
type Message struct {
	Id          string `json:"Id"`
	SenderID    string `json:"SenderID"`
	SenderName  string `json:"SenderName,omitempty"`
	SenderApodo string `json:"SenderApodo,omitempty"`
	GroupID     string `json:"GroupID"`
	Content     string `json:"Content"`
	Fecha       string `json:"Fecha"`
	AnswerId    string `json:"AnswerId"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*websocket.Conn),
		userGroups: make(map[string]map[string]bool),
		register:   make(chan RegisterClient),
		unregister: make(chan string),
		broadcast:  make(chan Message),
		aiChannel:  make(chan Message),
		done:       make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.done:
			log.Println("Hub: Deteniendo el hub...")
			return
		case reg := <-h.register:
			h.mu.Lock()
			h.clients[reg.UserID] = reg.Conn
			h.mu.Unlock()
			log.Printf("Hub: Usuario %s registrado.\n", reg.UserID)

		case userID := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[userID]; ok {
				delete(h.clients, userID)
				log.Printf("Hub: Usuario %s desconectado.\n", userID)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.Lock()
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Hub: Error al serializar el mensaje a JSON: %v", err)
				h.mu.Unlock()
				continue
			}

			for userID, conn := range h.clients {
				// Verificamos si este usuario pertenece al grupo del mensaje
				if h.userGroups[userID] != nil && h.userGroups[userID][msg.GroupID] {
					if err := conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
						log.Printf("Hub: Error al enviar a %s: %v", userID, err)
						h.unregister <- userID
					}
				}
			}
			h.aiChannel <- msg
			h.mu.Unlock()
		}
	}
}

// Register registra una conexión con su UserID
func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.register <- RegisterClient{UserID: userID, Conn: conn}
}

// Unregister desregistra una conexión
func (h *Hub) Unregister(userID string) {
	h.unregister <- userID
}

// Broadcast envía un mensaje a un grupo específico
func (h *Hub) Broadcast(msg Message) {
	h.broadcast <- msg
}

// Nueva función pública para acceder al canal de la IA
func (h *Hub) AIChannel() chan Message {
	return h.aiChannel
}

// SubscribeUserToGroups suscribe a un usuario a múltiples grupos
func (h *Hub) SubscribeUserToGroups(userID string, groupIDs []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.userGroups[userID] = make(map[string]bool)
	for _, groupID := range groupIDs {
		h.userGroups[userID][groupID] = true
		log.Printf("Usuario %s suscrito al grupo %s.\n", userID, groupID)
	}
}

// CheckUserInGroup verifica si un usuario pertenece a un grupo
func (h *Hub) CheckUserInGroup(userID, groupID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Verificamos si el usuario existe en el mapa de userGroups
	userMap, ok := h.userGroups[userID]
	if !ok {
		return false // El usuario no tiene grupos registrados
	}

	// Verificamos si el grupo existe en el mapa de grupos del usuario
	_, ok = userMap[groupID]
	return ok
}

// Shutdown cierra el canal done para detener el hub
func (h *Hub) Shutdown() {
	close(h.done)
}
