package ollama_dto

import "time"

type RequestEmbedAI struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type ResponseEmbedAI struct {
	Embeddings [][]float32 `json:"embeddings"`
}

type RequestChatAI struct {
	Model    string          `json:"model"`
	Messages []MessageChatAI `json:"messages"`
	Stream   bool            `json:"stream"`
}

type MessageChatAI struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseChatAI struct {
	Model              string        `json:"model"`
	CreatedAt          time.Time     `json:"created_at"`
	Message            MessageChatAI `json:"message"`
	Done               bool          `json:"done"`
	TotalDuration      int64         `json:"total_duration"`
	LoadDuration       int           `json:"load_duration"`
	PromptEvalCount    int           `json:"prompt_eval_count"`
	PromptEvalDuration int           `json:"prompt_eval_duration"`
	EvalCount          int           `json:"eval_count"`
	EvalDuration       int64         `json:"eval_duration"`
}
