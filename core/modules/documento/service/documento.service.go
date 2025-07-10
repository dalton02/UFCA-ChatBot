package documento_service

import (
	"fmt"
	documento_dto "licor_model/core/modules/documento/dto"
	ollama_service "licor_model/core/modules/ollama/service"
	"licor_model/core/server/shared"
	"strconv"
	"strings"
)

func BuscarDocumentoPorContextoDB(contexto string) (documento_dto.DocumentoDto, error) {
	var documento documento_dto.DocumentoDto
	query := `SELECT contexto,conteudo,link,id FROM documentos WHERE contexto=$1`
	result := shared.DB.QueryRow(query, contexto)
	err := result.Scan(&documento.Contexto, &documento.Conteudo, &documento.Link, &documento.Id)
	if err != nil {
		return documento, err
	}
	return documento, nil
}

func BuscarDocumentoPorSimilaridadeDB(conteudo string) ([]documento_dto.DocumentoDto, error) {
	var documentos []documento_dto.DocumentoDto
	vetor, err := ollama_service.GerarEmbedding(conteudo)
	if err != nil {
		return documentos, err
	}
	value := builderQueryVetor(vetor)
	query := `SELECT contexto,conteudo,link FROM documentos ORDER BY embedding <-> ` + value + ` LIMIT 5`
	result, err := shared.DB.Query(query)
	if err != nil {
		return documentos, err
	}
	for result.Next() {
		var doc documento_dto.DocumentoDto
		result.Scan(&doc.Contexto, &doc.Conteudo, &doc.Link)
		documentos = append(documentos, doc)
	}
	return documentos, nil
}

func SalvarEmbeddingDb(contexto string, conteudo string, link string, vetor [][]float32) error {
	vetorSQL := builderQueryVetor(vetor)

	documento, _ := BuscarDocumentoPorContextoDB(contexto)
	if documento.Id > 0 {
		fmt.Println("Dado atualizado")
		queryUpdate := `UPDATE documentos SET contexto = $1, conteudo = $2, link = $3, embedding = ` + vetorSQL + ` WHERE id = $4;`
		_, err := shared.DB.Exec(queryUpdate, contexto, conteudo, link, documento.Id)
		return err
	}
	fmt.Println("Novo dado inserido no banco")

	queryInsert := `INSERT INTO documentos (contexto,conteudo,link,embedding) VALUES ($1,$2,$3,` + vetorSQL + `);`
	_, err := shared.DB.Exec(queryInsert, contexto, conteudo, link)
	return err
}

func builderQueryVetor(vetor [][]float32) string {
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
