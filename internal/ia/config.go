package ia

import "os"

type IAConfig struct {
	UserID     string
	LLMBaseURL string
	LLMName    string
	LLMAPIKey  string
	IsPromt    bool
}

var AiConfigurations = []IAConfig{
	{
		UserID:     "5",
		LLMBaseURL: "http://localhost:11434/v1",
		LLMName:    "gpt-oss:120b",
		LLMAPIKey:  os.Getenv("LLM_API_KEY_1"),
	},
	{
		UserID:     "6",
		LLMBaseURL: "http://localhost:1234/v1/chat/completions",
		LLMName:    "google/gemma-3-27b",
		LLMAPIKey:  os.Getenv("LLM_API_KEY_2"),
		IsPromt:    false,
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
