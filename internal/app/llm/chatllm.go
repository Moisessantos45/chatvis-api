package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Define las estructuras para el cuerpo de la solicitud y la respuesta
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatResponse struct {
	Model      string      `json:"model"`
	RemoteHost string      `json:"remote_host"`
	CreatedAt  string      `json:"created_at"`
	Message    ChatMessage `json:"message"`
	Done       bool        `json:"done"`
}
type CompletionBody struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	Stream      bool          `json:"stream"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

type OllamaCompletionBody struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"` // <--- Nuevo campo
	Temperature float64 `json:"temperature"`
	Stream      bool    `json:"stream"`
}

// La respuesta de Ollama Completion es diferente a la de Chat.
// Un ejemplo simple podría ser:
type OllamaCompletionResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"` // <--- Campo de contenido
	CreatedAt string `json:"created_at"`
	Done      bool   `json:"done"`
}

// func PostCompletion(messages []ChatMessage, baseURL string, nameLLM string, apiKey string) (string, error) {
// 	requestBody := CompletionBody{
// 		Model:       nameLLM,
// 		Messages:    messages,
// 		Temperature: 0.7,
// 		Stream:      false,
// 	}

// 	jsonBody, err := json.Marshal(requestBody)
// 	if err != nil {
// 		return "", fmt.Errorf("fallo al serializar el cuerpo de la solicitud: %w", err)
// 	}

// 	client := &http.Client{Timeout: 30 * time.Second}

// 	req, err := http.NewRequest("POST", "https://n8n.glimpse.uaslp.mx/ollama/api/generate"), bytes.NewBuffer(jsonBody))
// 	if err != nil {
// 		return "", fmt.Errorf("fallo al crear la solicitud: %w", err)
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("fallo al enviar la solicitud: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	bodyBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("fallo al leer la respuesta: %w", err)
// 	}

// 	fmt.Println(">>> Respuesta cruda del LLM:\n", string(bodyBytes))

// 	// 1. Intentar decodificar como OpenAI
// 	var completionResponse ChatCompletionResponse
// 	if err := json.Unmarshal(bodyBytes, &completionResponse); err == nil && len(completionResponse.Choices) > 0 {
// 		return completionResponse.Choices[0].Message.Content, nil
// 	}

// 	// 2. Intentar decodificar como Ollama
// 	var ollamaResponse OllamaChatResponse
// 	if err := json.Unmarshal(bodyBytes, &ollamaResponse); err == nil && ollamaResponse.Message.Content != "" {
// 		return ollamaResponse.Message.Content, nil
// 	}

// 	return "", fmt.Errorf("no se encontró contenido válido en la respuesta del LLM")
// }

func PostCompletionOllamaPrompt(messages []ChatMessage, baseURL string, nameLLM string, apiKey string) (string, error) {
	// Asegurarse de que haya al menos un mensaje para el prompt
	if len(messages) == 0 {
		return "", fmt.Errorf("se requiere al menos un mensaje para el prompt")
	}
	log.Println("PostCompletionOllamaPrompt - Mensajes recibidos:", messages)
	// El 'prompt' para el endpoint /api/generate suele ser el contenido del último mensaje (el del usuario).
	// Para un chat, podrías querer concatenar todos los mensajes. Asumiremos el último mensaje por simplicidad.
	prompt := messages[len(messages)-1].Content

	requestBody := OllamaCompletionBody{
		Model:       nameLLM,
		Prompt:      prompt, // Usa el campo 'prompt'
		Temperature: 0.7,
		Stream:      false, // Como en tu curl
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("fallo al serializar el cuerpo de la solicitud: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	// Nota: El curl usa "https://n8n.glimpse.uaslp.mx/ollama/api/generate"
	// Usaremos la baseURL y el path "/api/generate"
	// url := fmt.Sprintf("%s/api/generate")

	req, err := http.NewRequest("POST", "https://n8n.glimpse.uaslp.mx/ollama/api/generate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("fallo al crear la solicitud: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	// Tu curl usa 'Authorization: Bearer <token>', así que replicamos eso:
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fallo al enviar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("fallo al leer la respuesta: %w", err)
	}

	fmt.Println(">>> Respuesta cruda del LLM (Ollama Completion):\n", string(bodyBytes))

	// Intentar decodificar como Ollama Completion
	var ollamaCompletion OllamaCompletionResponse
	if err := json.Unmarshal(bodyBytes, &ollamaCompletion); err == nil {
		if ollamaCompletion.Response != "" {
			return ollamaCompletion.Response, nil
		}
	}

	return "", fmt.Errorf("no se encontró contenido válido en la respuesta del LLM (Ollama Completion)")
}
