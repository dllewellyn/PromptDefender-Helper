package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func main() {
	r := gin.Default()

	// Load HTML template
	r.SetHTMLTemplate(template.Must(template.New("hello").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Hello, World!</title>
		</head>
		<body>
			<h1>Hello and welcome, {{.}}</h1>
		</body>
		</html>
	`)))

	// Define a route and handler
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello", "gopher")
	})

	// Start the server
	r.Run(":8080")
}
