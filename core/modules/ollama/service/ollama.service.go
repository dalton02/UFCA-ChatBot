package ollama_service

import (
	"bytes"
	"encoding/json"
	ollama_dto "licor_model/core/modules/ollama/dto"
	"net/http"
	"os"
)

func SendRequest(ollamaReq ollama_dto.Request) (*ollama_dto.Response, error) {
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
	ollamaResp := ollama_dto.Response{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	return &ollamaResp, err
}
