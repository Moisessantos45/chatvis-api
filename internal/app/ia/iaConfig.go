package ia

// ia_config.go (Ejemplo de un nuevo modelo)
type IAConfig struct {
	UserID     string
	LLMBaseURL string
	LLMName    string
	LLMAPIKey  string
	IsPromt   bool
}

var AiConfigurations = []IAConfig{
	{
		UserID:     "14",
		LLMBaseURL: "https://n8n.glimpse.uaslp.mx/ollama/api/generate",
		LLMName:    "llama3.2",
		LLMAPIKey:  "bc8af4b4-b264-4fae-b748-324693ab0151",
		IsPromt: false,
	},
	{
		UserID:     "15",
		LLMBaseURL: "http://localhost:11434",
		LLMName:    "gpt-oss:latest",
		LLMAPIKey:  "bc8af4b4-b264-4fae-b748-324693ab0151",
		IsPromt: false,
	},
	{
		UserID:     "16",
		LLMBaseURL: "http://localhost:11434",
		LLMName:    "gemma3:27b",
		LLMAPIKey:  "bc8af4b4-b264-4fae-b748-324693ab0151",
		IsPromt: false,
	},
	{
		UserID:     "17",
		LLMBaseURL: "http://localhost:11434",
		LLMName:    "qwen2.5:32b",
		LLMAPIKey:  "bc8af4b4-b264-4fae-b748-324693ab0151",
		IsPromt: false,
	},
}
