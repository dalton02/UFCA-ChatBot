package util

import (
	"fmt"
	"io"
	"licor_model/core/util/executor"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
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

func ConvertFloat32SliceToString(data [][]float32) string {
	var outerBuilder strings.Builder
	outerBuilder.WriteString("[") // Começa com o colchete externo

	for i, innerSlice := range data {
		if i > 0 {
			outerBuilder.WriteString(",") // Adiciona vírgula entre os slices internos
		}
		outerBuilder.WriteString("[") // Começa com o colchete do slice interno

		for j, val := range innerSlice {
			if j > 0 {
				outerBuilder.WriteString(",") // Adiciona vírgula entre os valores
			}
			// Converte float32 para string e adiciona aspas duplas
			outerBuilder.WriteString(`"` + strconv.FormatFloat(float64(val), 'f', -1, 32) + `"`)
		}
		outerBuilder.WriteString("]") // Fecha o colchete do slice interno
	}
	outerBuilder.WriteString("]") // Fecha o colchete externo

	return outerBuilder.String()
}

func BuilderQueryVetor(vetor [][]float32) string {
	builderVetores := strings.Builder{}
	for i, innerSlice := range vetor {
		builderVetores.WriteString("(")
		builderVetores.WriteString(`'[`)
		for j, val := range innerSlice {
			builderVetores.WriteString(strconv.FormatFloat(float64(val), 'f', -1, 32))
			if j < len(innerSlice)-1 {
				builderVetores.WriteString(",")
			}
		}

		builderVetores.WriteString(`]'`)
		builderVetores.WriteString(")")
		if i < len(vetor)-1 {
			builderVetores.WriteString(",")
		}
	}
	return builderVetores.String()
}

func CountSQL(executor executor.Executor, expression string, builder *goqu.SelectDataset) (total int, err error) {
	countQuery := builder.
		ClearSelect().
		ClearOrder().
		ClearOffset().
		ClearLimit()

	countQuery = countQuery.Select(goqu.COUNT(goqu.L(expression)))

	sql, args, err := countQuery.ToSQL()

	if err != nil {
		return 0, fmt.Errorf("erro ao gerar SQL: %w", err)
	}

	rows, err := executor.Query(sql, args...)

	if err != nil {
		return 0, fmt.Errorf("erro ao executar COUNT: %w", err)
	}

	for rows.Next() {
		var tmp int
		rows.Scan(&tmp)
		total += tmp
	}

	return total, nil
}
