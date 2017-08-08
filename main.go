package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/heroku/x/hmetrics"
)

type fataler interface {
	Fatal() bool
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	hml := log.New(os.Stderr, "heroku metrics", 0)
	if err := hmetrics.Report(context.Background(), hml); err != nil {
		if f, ok := err.(fataler); ok {
			if f.Fatal() {
				log.Fatal(err)
			}
			log.Println(err)
		}
	}

	router.Run(":" + port)
}
