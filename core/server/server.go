package server

import (
	"fmt"
	"licor_model/core/server/shared"
	"licor_model/core/util/interceptor"
	"licor_model/docs"
	"os"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RoutesH struct {
	Engine *gin.Engine
	Groups RouteGroups
}

type RouteGroups struct {
	PublicGroup *gin.RouterGroup //Agrupamento sem middleware
	JwtGroup    *gin.RouterGroup //Agrupamento que passa por validação jwt
	OfficeGroup *gin.RouterGroup //Agrupamento que passa por validação de escritorio e autoridade do usuário
}

var routes *RoutesH

func MainServer() {

	StartEngine()

	InitInjections()

	HandleDocs()

	routes.Engine.Run(":3000")

}

func StartEngine() {
	interceptor.InitValidator()

	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	routes = &RoutesH{
		Engine: gin.Default(),
	}

	routes.Engine.Use(shared.Cors)
	routes.Engine.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	routes.Engine.Use(gin.Logger())
	routes.Engine.Use(gin.Recovery())

}

func HandleDocs() {

	docs.SwaggerInfo.Host = os.Getenv("BACKEND_HOST")
	routes.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.Engine.GET("/docs", func(ctx *gin.Context) {

		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",

			CustomOptions: scalar.CustomOptions{
				PageTitle: "CHATBOT - DOCUMENTAÇÃO",
			},
			DarkMode: true,
			Layout:   scalar.LayoutModern,
			Theme:    scalar.ThemePurple,

			Authentication: "jwt",
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(ctx.Writer, htmlContent)

	})
}
