package server

import (
	example_controller "licor_model/core/modules/example.module/controller"
	example_dto "licor_model/core/modules/example.module/dto"

	"github.com/dalton02/licor/licor"
)

func MainServer() {

	licor.Public[example_dto.RouteTestValidation, any]("/exemple").Get(example_controller.ExampleRoute)
	licor.Init("4000")

}
