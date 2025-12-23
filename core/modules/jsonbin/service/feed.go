package jsonbin_service

import (
	"fmt"
	document_dto "licor_model/core/modules/document/dto"
	jsonbin_dto "licor_model/core/modules/jsonbin/dto"
	"sync"
)

// Alimenta a IA a partir de um dado arquivo em uma dada época
func (service *JsonBinService) FeedFromFile(year int, month int, day int, filename string) error {

	path := service.makePathWithDate(year, month, day, filename)
	data, err := service.readJsonFile(path)
	if err != nil {
		return err
	}

	err = service.documentService.DeleteAllFromOrigin(document_dto.JsonStudents)
	if err != nil {
		return err
	}
	for _, d := range data {

		err = service.documentService.UpsertDocument(d.Context, d.Content, "", document_dto.JsonStudents)
		if err != nil {
			fmt.Println(err.Error())
		}

	}

	return nil
}

// Alimenta a IA recebendo um array de json
func (service *JsonBinService) FeedIAWithData(data []jsonbin_dto.JsonFile) error {

	err := service.documentService.DeleteAllFromOrigin(document_dto.JsonStudents)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, doc := range data {

		wg.Add(1)
		go func() {
			service.documentService.UpsertDocument(doc.Context, doc.Content, "", document_dto.JsonStudents)
			wg.Done()
		}()

	}

	wg.Wait()

	return nil
}
