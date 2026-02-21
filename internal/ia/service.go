package ia

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/llm"
	"chatvis-chat/internal/websocket"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type AIService struct {
	Hub            *websocket.Hub
	MensajeUseCase domain.MensajeUseCase
	GrupoUseCase   domain.GrupoUseCase
	conversations  sync.Map

	Config IAConfig

	inputChannel chan websocket.Message
	jobs         chan websocket.Message
	quit         chan struct{}
	stopOnce     sync.Once
}

func NewAIService(h *websocket.Hub, mr domain.MensajeUseCase, gu domain.GrupoUseCase, config IAConfig) *AIService {
	return &AIService{
		Hub:            h,
		MensajeUseCase: mr,
		GrupoUseCase:   gu,
		quit:           make(chan struct{}),
		Config:         config,
		inputChannel:   make(chan websocket.Message, 100), // Buffer para manejar ráfagas de mensajes
		jobs:           make(chan websocket.Message, 100), // Canal de trabajos para el pool de workers
	}
}

func (s *AIService) InputChannel() chan websocket.Message {
	return s.inputChannel
}

// Start listens for new messages from the Hub and decides whether to respond.
func (s *AIService) Start(ctx context.Context, workerCount int) {
	// Primero suscribir a los grupos
	if err := s.SuscribeToGroup(); err != nil {
		log.Printf("Error al suscribir IA a grupos: %v", err)
	} else {
		log.Printf("IA suscrita a sus grupos correctamente")
	}

	for i := 0; i < workerCount; i++ {
		go s.worker(ctx, i)
	}

	go s.listenForMessages(ctx)
}

func (s *AIService) SuscribeToGroup() error {
	idUsuario, err := strconv.ParseUint(s.Config.UserID, 10, 64)
	if err != nil {
		return err
	}

	aiGroupClaves, err := s.GrupoUseCase.GetAllGruposByUsuarioIdToClaves(idUsuario)
	if err != nil {
		return err
	}

	s.Hub.SubscribeUserToGroups(s.Config.UserID, aiGroupClaves)

	return nil
}

func (s *AIService) listenForMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("AIService: Contexto cancelado, deteniendo gorutina de escucha.")
			return

		case msg, ok := <-s.inputChannel:
			if !ok {
				log.Println("AIService: Canal de entrada cerrado, deteniendo gorutina de escucha.")
				return
			}

			log.Printf("AIService: recibido mensaje en AIChannel: %+v", msg)

			// 1. Opcional: Evita que la IA responda a sus propios mensajes.
			if msg.SenderID == s.Config.UserID {
				continue
			}

			// 2. Verifica si la IA está en el grupo.
			if !s.Hub.CheckUserInGroup(s.Config.UserID, msg.GroupID) {
				log.Printf("AIService: La IA no está en el grupo %s, no responde.", msg.GroupID)
				continue
			}

			// Mandar al pool de workers
			select {
			case <-ctx.Done():
				return
			case s.jobs <- msg:
			}
		}
	}
}

func (s *AIService) worker(ctx context.Context, id int) {
	log.Printf("Worker %d iniciado", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d detenido", id)
			return
		case msg, ok := <-s.jobs:
			if !ok {
				log.Printf("Worker %d: canal jobs cerrado", id)
				return
			}
			log.Printf("Worker %d procesando mensaje %+v", id, msg)
			s.generateAndSendResponse(msg)
		}
	}
}

func (s *AIService) generateAndSendResponse(incomingMsg websocket.Message) {

	aiUserID, err := strconv.ParseUint(s.Config.UserID, 10, 64)
	if err != nil {
		log.Printf("Error al convertir AI User ID: %v", err)
		return
	}

	// 3.1 Obtener ID del Grupo
	responseGrupo, err := s.GrupoUseCase.GetByClave(incomingMsg.GroupID)
	if err != nil || responseGrupo == nil {
		log.Printf("Error al obtener el grupo por clave %s: %v", incomingMsg.GroupID, err)
		return
	}
	grupoIDUint := responseGrupo.Id

	// 3.2. Obtener los mensajes nuevos del grupo de la base de datos
	allGroupMessages, err := s.MensajeUseCase.GetNuevosMensajesParaIA(aiUserID, grupoIDUint)
	if err != nil {
		log.Printf("Error al obtener mensajes nuevos del grupo %s: %v", incomingMsg.GroupID, err)
		return
	}

	if len(allGroupMessages) == 0 {
		log.Println("No hay mensajes nuevos para procesar por la IA.")
		return
	}

	lastMessageProcessed := allGroupMessages[len(allGroupMessages)-1].Id

	// 3.3. Convertir los mensajes a un formato que el LLM entienda
	llmMessages := s.buildPromptFromHistory(allGroupMessages)

	// 3.4. Llamar a la función del cliente LLM con todo el historial
	aiResponse, err := llm.PostCompletionOllamaPrompt(llmMessages, s.Config.LLMBaseURL, s.Config.LLMName, s.Config.LLMAPIKey)
	if err != nil {
		log.Printf("Error al generar respuesta de IA: %v", err)
		return
	}

	if aiResponse == "" {
		log.Println("Respuesta de IA vacía, no se envía el mensaje.")
		return
	}

	// 3.5. Crear y enviar el mensaje de la IA al Hub tenemos que usar la funcion ParseAIResponse
	aiMsg, err := s.ParseAIResponse(aiResponse, incomingMsg.GroupID, s.Config.UserID)
	if err != nil {
		log.Printf("Error al parsear la respuesta de IA: %v", err)
		return
	}

	// 3.6. Guardar el mensaje de la IA en la base de datos
	gormMsg, err := s.saveAIToDB(aiMsg)
	if err != nil {
		log.Printf("Error al guardar el mensaje de IA en la base de datos: %v", err)
		return
	}

	// Actualizar el ID del mensaje en el objeto aiMsg para que el Hub tenga el ID correcto
	aiMsg.Id = strconv.FormatUint(gormMsg.Id, 10)

	// Actualizar checkpoint (el ultimo mensaje que leyó la IA)
	err = s.MensajeUseCase.ActualizarPuntoControl(aiUserID, grupoIDUint, lastMessageProcessed)
	if err != nil {
		log.Printf("Error al actualizar punto de control de IA: %v", err)
	}

	// 3.7. Enviar el mensaje de la IA a través del Hub
	s.Hub.Broadcast(*aiMsg)
}

func (s *AIService) buildPromptFromHistory(mensajes []domain.Mensaje) []llm.ChatMessage {
	var chatMessages []llm.ChatMessage
	aiUserIDUint, _ := strconv.ParseUint(s.Config.UserID, 10, 64)
	if s.Config.IsPromt {
		chatMessages = append(chatMessages, llm.ChatMessage{
			Role:    "system",
			Content: s.getSystemPrompt()})
		aiUserIDUint, _ = strconv.ParseUint(s.Config.UserID, 10, 64)
	}
	// Itera sobre los mensajes de la base de datos para construir el historial.
	for _, msg := range mensajes {
		role := "user"
		if msg.UsuarioId == aiUserIDUint && s.Config.IsPromt {
			role = "assistant"
		}

		chatMessages = append(chatMessages, llm.ChatMessage{
			Role:    role,
			Content: msg.Contenido,
		})
	}
	return chatMessages
}

func (s *AIService) ParseAIResponse(aiResponseJSON string, groupID string, senderID string) (*websocket.Message, error) {
	// Usar interface{} para manejar tanto números como strings
	var parsedResponse struct {
		AnswerID interface{} `json:"answer_id"`
		Content  string      `json:"content"`
	}

	err := json.Unmarshal([]byte(aiResponseJSON), &parsedResponse)
	if err != nil {
		return nil, err
	}

	var answerID string

	// Manejar diferentes tipos de answer_id
	switch v := parsedResponse.AnswerID.(type) {
	case nil:
		// null -> cadena vacía
		answerID = ""
	case float64:
		// JSON números se parsean como float64
		answerID = strconv.FormatUint(uint64(v), 10)
	case string:
		// Si viene como string, validar que sea numérico
		if v != "" {
			if _, err := strconv.ParseUint(v, 10, 64); err != nil {
				log.Printf("ADVERTENCIA: answer_id string inválido '%s', usando null", v)
				answerID = ""
			} else {
				answerID = v
			}
		} else {
			answerID = ""
		}
	default:
		log.Printf("ADVERTENCIA: answer_id tipo desconocido %T, usando null", v)
		answerID = ""
	}

	aiMsg := &websocket.Message{
		SenderID: senderID,
		GroupID:  groupID,
		Content:  parsedResponse.Content,
		Fecha:    time.Now().Format(time.RFC3339),
		AnswerId: answerID,
	}

	return aiMsg, nil
}

func (s *AIService) saveAIToDB(aiMsgHub *websocket.Message) (*domain.Mensaje, error) {
	responseGrupo, err := s.GrupoUseCase.GetByClave(aiMsgHub.GroupID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el grupo por clave: %w", err)
	}

	if responseGrupo == nil {
		return nil, fmt.Errorf("grupo no encontrado para la clave: %s", aiMsgHub.GroupID)
	}

	groupID := responseGrupo.Id

	senderID, err := strconv.ParseUint(aiMsgHub.SenderID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error al convertir SenderID: %w", err)
	}

	var responseID *uint64
	// Manejar "-1" como caso especial (comentario general)
	if aiMsgHub.AnswerId != "" && aiMsgHub.AnswerId != "-1" {
		parsedResponseID, err := strconv.ParseUint(aiMsgHub.AnswerId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error al convertir AnswerId: %w", err)
		}
		responseID = &parsedResponseID
	}

	gormMsg := &domain.Mensaje{
		Contenido:  aiMsgHub.Content,
		Fecha:      time.Now(),
		GrupoId:    groupID,
		UsuarioId:  senderID,
		ResponseId: responseID,
	}

	response, err := s.MensajeUseCase.Create(gormMsg)
	if err != nil {
		return nil, fmt.Errorf("error al guardar el mensaje de la IA: %w", err)
	}

	gormMsg.Id = response.Id
	return gormMsg, nil
}

func (s *AIService) getSystemPrompt() string {
	return `FORMATO JSON OBLIGATORIO:
Debes responder SOLO con JSON válido. NO agregues texto adicional antes o después del JSON.

CRÍTICO - REGLAS DE answer_id:
- Si respondes a un mensaje específico: usa el ID numérico (ejemplo: 223, 45, 1) SIN comillas
- Si es un comentario general que NO responde a ningún mensaje: usa null (sin comillas)
- answer_id NUNCA puede ser:
  * Una palabra o string (como "ayuda", "comentario", "general", "-1")
  * Cadena vacía ""
  * Números entre comillas como "223" o "-1"

EJEMPLOS CORRECTOS:

Respondiendo al mensaje con ID 223:
{
  "answer_id": 223,
  "content": "Dale, cual es?"
}

Comentario general (no responde a nadie en específico):
{
  "answer_id": null,
  "content": "yo puedo con eso"
}

EJEMPLOS INCORRECTOS - NUNCA HAGAS ESTO:
{
  "answer_id": "ayuda",
  "content": "..."
}

{
  "answer_id": "223",
  "content": "..."
}

{
  "answer_id": "-1",
  "content": "..."
}

RECUERDA: answer_id es un número sin comillas (223) o null. Si no sabes a qué mensaje responder, usa null.`
}

// Stop detiene la gorutina que escucha los mensajes.
func (s *AIService) Stop() {
	s.stopOnce.Do(func() {
		close(s.quit)
		close(s.inputChannel)
		close(s.jobs)
	})
}
