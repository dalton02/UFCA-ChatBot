package jsonbin_dto

type JsonFile struct {
	Context string `json:"context"`
	Content string `json:"content"`
}

type JsonBinTree map[int]map[int]map[int][]string
