package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo/v4"
)

func main() {
	stage := os.Getenv("AWS_STAGE")
	if stage == "local" {
		e := echo.New()
		e.POST("/client", GetIndex)
		e.Logger.Fatal(e.Start(":1323"))
	} else {
		lambda.Start(Handler)
	}
}
