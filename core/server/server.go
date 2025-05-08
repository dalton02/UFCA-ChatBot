package server

import (
	example_controller "licor_model/core/modules/example.module/controller"
	_ "licor_model/docs"

	"github.com/dalton02/licor/licor"
)

func MainServer() {
	licor.Public[any, any]("/hello-world").Get(example_controller.ExampleRoute)

	licor.Init("4000")

}
