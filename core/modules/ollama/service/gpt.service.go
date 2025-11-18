package ai_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	ollama_dto "licor_model/core/modules/ollama/dto"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func init() {
	_ = godotenv.Load()
}

func SendRequestOpenAIStream(c *gin.Context, links []string, request ollama_dto.RequestChatAI) (*ollama_dto.ResponseChatAI, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY não encontrada")
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	client := openai.NewClient(apiKey)
	ctx := c.Request.Context()

	var messages []openai.ChatCompletionMessage
	for _, m := range request.Messages {
		role := openai.ChatMessageRoleUser
		if m.Role == "assistant" {
			role = openai.ChatMessageRoleAssistant
		} else if m.Role == "system" {
			role = openai.ChatMessageRoleSystem
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: m.Content,
		})
	}

	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.7,
		Stream:      true,
	})
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	var fullContent strings.Builder

	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			c.Writer.Write([]byte("data: [DONE]\n\n"))
			c.Writer.Flush()
			break
		}

		if err != nil {
			errorMsg := fmt.Sprintf("data: {\"error\": \"%s\"}\n\n", err.Error())
			c.Writer.Write([]byte(errorMsg))
			c.Writer.Flush()
			return nil, err
		}

		if len(response.Choices) == 0 {
			continue
		}

		delta := response.Choices[0].Delta
		index := response.Choices[0].Index

		if delta.Content != "" {
			content := delta.Content
			fullContent.WriteString(content)

			payload := map[string]any{
				"content": content,
				"index":   index,
				"links":   links,
				"model":   response.Model,
			}
			if delta.Role != "" {
				payload["role"] = delta.Role
			}

			jsonData, _ := json.Marshal(payload)
			fmt.Fprintf(c.Writer, "data: %s\n\n", jsonData)
			c.Writer.Flush()
		}

	}

	finalResponse := &ollama_dto.ResponseChatAI{
		Model: model,
		Message: ollama_dto.MessageChatAI{
			Role:    "assistant",
			Content: fullContent.String(),
		},
		CreatedAt: time.Now(),
		Done:      true,
	}

	return finalResponse, nil
}
