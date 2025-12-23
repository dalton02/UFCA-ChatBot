package jsonbin_service

import (
	"encoding/json"
	"fmt"
	jsonbin_dto "licor_model/core/modules/jsonbin/dto"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func (service *JsonBinService) makeDatePath(userEmail string) string {
	currentDate := time.Now()

	year := currentDate.Year()
	month := int(currentDate.Month())
	day := currentDate.Day()
	hour := currentDate.Hour()

	fileName := fmt.Sprintf(
		"%s_%02d.json",
		service.userKey(userEmail),
		hour,
	)

	return filepath.Join(
		service.BasePath,
		fmt.Sprintf("%04d", year),
		fmt.Sprintf("%02d", month),
		fmt.Sprintf("%02d", day),
		fileName,
	)
}

func (service *JsonBinService) makePathWithDate(year int, month int, day int, filename string) string {

	fileName := fmt.Sprintf(
		"%s",
		filename,
	)

	return filepath.Join(
		service.BasePath,
		fmt.Sprintf("%04d", year),
		fmt.Sprintf("%02d", month),
		fmt.Sprintf("%02d", day),
		fileName,
	)
}

func (service *JsonBinService) readJsonFile(path string) ([]jsonbin_dto.JsonFile, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data []jsonbin_dto.JsonFile
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (service *JsonBinService) saveJsonFile(path string, data []jsonbin_dto.JsonFile) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}

func (service *JsonBinService) userKey(email string) string {
	encoded := url.QueryEscape(email)
	return encoded
}
