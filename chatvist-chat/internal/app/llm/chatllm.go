package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func PostCompletion(messages []ChatMessage, baseURL string, nameLLM string, apiKey string) (string, error) {
	requestBody := CompletionBody{
		Model:       nameLLM,
		Messages:    messages,
		Temperature: 0.7,
		Stream:      false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("fallo al serializar el cuerpo de la solicitud: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s", baseURL), bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("fallo al crear la solicitud: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
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

	fmt.Println(">>> Respuesta cruda del LLM:\n", string(bodyBytes))

	// 1. Intentar decodificar como OpenAI
	var completionResponse ChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &completionResponse); err == nil && len(completionResponse.Choices) > 0 {
		return completionResponse.Choices[0].Message.Content, nil
	}

	// 2. Intentar decodificar como Ollama
	var ollamaResponse OllamaChatResponse
	if err := json.Unmarshal(bodyBytes, &ollamaResponse); err == nil && ollamaResponse.Message.Content != "" {
		return ollamaResponse.Message.Content, nil
	}

	return "", fmt.Errorf("no se encontró contenido válido en la respuesta del LLM")
}
