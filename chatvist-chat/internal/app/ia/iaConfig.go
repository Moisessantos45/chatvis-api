package ia

import "os"

// ia_config.go (Ejemplo de un nuevo modelo)
type IAConfig struct {
	UserID     string
	LLMBaseURL string
	LLMName    string
	LLMAPIKey  string
	Isprompt   bool
}

var AiConfigurations = []IAConfig{
	// {
	// 	UserID:     "5",
	// 	LLMBaseURL: "http://localhost:1234/api/chat",
	// 	LLMName:    "gpt-oss:120b ",
	// 	LLMAPIKey:  os.Getenv("LLM_API_KEY_1"),
	// },
	{
		UserID:     "6",
		LLMBaseURL: "http://localhost:1234/v1/chat/completions",
		LLMName:    "google/gemma-3-27b",
		LLMAPIKey:  os.Getenv("LLM_API_KEY_2"),
		Isprompt: false,
	},
	// {
	// 	UserID:     "7",
	// 	LLMBaseURL: "http://localhost:1234/api/chat",
	// 	LLMName:    "gemma3:27b",
	// 	LLMAPIKey:  os.Getenv("LLM_API_KEY_3"),
	// },
	// {
	// 	UserID:     "8",
	// 	LLMBaseURL: "http://localhost:1234/api/chat",
	// 	LLMName:    "qwen3-coder:480b-cloud",
	// 	LLMAPIKey:  os.Getenv("LLM_API_KEY_4"),
	// },
}
