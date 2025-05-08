package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func NormalizeString(input string) (string, error) {
	return url.QueryUnescape(input)

}
func StringToDate(input string) time.Time {
	parts := strings.Split(input, ".")
	if len(parts) < 2 {
		return time.Now()
	}

	datePart := strings.TrimSpace(parts[0][12:])

	publishDate, err := time.Parse("02/01/2006", datePart)
	if err != nil {
		return time.Now()
	}
	return publishDate
}

func ArrayToInt(tmpString []string) (tmpInt []int) {
	var tmp []int

	for _, str := range tmpString {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil
		}
		tmp = append(tmp, num)
	}
	return tmp
}

func ArrayToString(array []int) string {
	// Converter cada número em uma string
	var strNumbers []string

	for _, num := range array {
		strNumbers = append(strNumbers, fmt.Sprintf("%d", num))
	}

	numbersString := strings.Join(strNumbers, ",")

	result := fmt.Sprintf("{%s}", numbersString)
	return result
}

func SanatizeFile(file multipart.File) bool {

	buffer := make([]byte, 512)
	_, err := file.Read(buffer) //Lendo 512 bytes iniciais do arquivos, necessario realocação
	if err != nil {
		return false
	}

	mime := fmt.Sprintf("%s", http.DetectContentType(buffer))
	if mime != "image/jpeg" && mime != "image/png" && mime != "image/webp" {
		return false
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return false
	}
	return true
}

func CommonToArray(data string) []string {
	var array []string
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			newString := data[start:i]
			array = append(array, newString)
			start = i + 1
		}
	}

	if start < len(data) {
		newString := data[start:]
		array = append(array, newString)
	}

	return array

}
