package jsonbin_service

import (
	"fmt"
	document_service "licor_model/core/modules/document/service"
	jsonbin_dto "licor_model/core/modules/jsonbin/dto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type JsonBinService struct {
	BasePath        string
	documentService *document_service.DocumentService
}

func NewJsonBinService(docService *document_service.DocumentService) *JsonBinService {
	return &JsonBinService{
		BasePath:        "bins/json",
		documentService: docService,
	}
}

func (service *JsonBinService) SaveContext(userEmail string, data []jsonbin_dto.JsonFile) error {
	filePath := service.makeDatePath(userEmail)
	return service.saveJsonFile(filePath, data)
}

func (service *JsonBinService) GetFromPath(year, month, day int, filename string) ([]jsonbin_dto.JsonFile, error) {

	path := service.makePathWithDate(year, month, day, filename)
	data, err := service.readJsonFile(path)
	if err != nil {
		return data, err
	}
	return data, nil

}

func (service *JsonBinService) GetLatestUserData(userEmail string) ([]jsonbin_dto.JsonFile, error) {
	user := service.userKey(userEmail)

	var latestPath string
	var latestTime time.Time

	err := filepath.Walk(service.BasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if !strings.HasPrefix(info.Name(), user+"_") {
			return nil
		}

		stat, err := os.Stat(path)
		if err != nil {
			return nil
		}

		if stat.ModTime().After(latestTime) {
			latestTime = stat.ModTime()
			latestPath = path
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if latestPath == "" {
		return nil, os.ErrNotExist
	}

	fmt.Println(latestPath)
	return service.readJsonFile(latestPath)
}
func (service *JsonBinService) GetLatestData() ([]jsonbin_dto.JsonFile, error) {
	var latestPath string
	var latestTime time.Time

	err := filepath.Walk(service.BasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		stat, err := os.Stat(path)
		if err != nil {
			return nil
		}

		if stat.ModTime().After(latestTime) {
			latestTime = stat.ModTime()
			latestPath = path
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if latestPath == "" {
		return nil, os.ErrNotExist
	}

	return service.readJsonFile(latestPath)
}

func (service *JsonBinService) GetTree() (jsonbin_dto.JsonBinTree, error) {
	tree := make(jsonbin_dto.JsonBinTree)

	err := filepath.WalkDir(service.BasePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) != ".json" {
			return nil
		}

		rel, err := filepath.Rel(service.BasePath, path)
		if err != nil {
			return nil
		}

		str := filepath.SplitList(rel)[0]
		parts := strings.Split(str, "/")

		fmt.Println(parts, len(parts))
		if len(parts) != 4 {
			return nil
		}

		year, err1 := strconv.Atoi(parts[0])
		month, err2 := strconv.Atoi(parts[1])
		day, err3 := strconv.Atoi(parts[2])

		fmt.Println(year, month, day)

		if err1 != nil || err2 != nil || err3 != nil {
			return nil
		}

		if _, ok := tree[year]; !ok {
			tree[year] = make(map[int]map[int][]string)
		}

		if _, ok := tree[year][month]; !ok {
			tree[year][month] = make(map[int][]string)
		}

		tree[year][month][day] = append(
			tree[year][month][day],
			parts[3],
		)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tree, nil
}
