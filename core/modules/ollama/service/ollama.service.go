package ollama_service

import (
	"bytes"
	"encoding/json"
	ollama_dto "licor_model/core/modules/ollama/dto"
	"net/http"
	"os"
)

func GerarEmbedding(conteudo string) ([][]float32, error) {
	var vetor [][]float32
	url := os.Getenv("OLLAMA_URL") + "/api/embed"
	content := ollama_dto.RequestEmbedAI{
		Model: os.Getenv("OLLAMA_EMBED_MODEL"),
		Input: conteudo,
	}
	js, _ := json.Marshal(&content)
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return vetor, err
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return vetor, err
	}
	defer httpResp.Body.Close()
	var result ollama_dto.ResponseEmbedAI
	json.NewDecoder(httpResp.Body).Decode(&result)
	vetor = result.Embeddings
	return vetor, err
}

func SendRequest(ollamaReq ollama_dto.RequestChatAI) (*ollama_dto.ResponseChatAI, error) {
	url := os.Getenv("OLLAMA_URL") + "/api/chat"
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	ollamaResp := ollama_dto.ResponseChatAI{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	return &ollamaResp, err
}
