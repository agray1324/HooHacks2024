package server

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func Router() *gin.Engine {
    r := gin.Default()
	r.LoadHTMLGlob("server/resources/*")
    r.GET("/", func(c *gin.Context) {
        c.Header("Content-Type", "text/html")
        c.HTML(http.StatusOK, "index.tmpl", gin.H{})
    })

    return r
}