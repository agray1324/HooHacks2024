package server

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func Router() *gin.Engine {
    r := gin.Default()
	r.LoadHTMLGlob("server/resources/*.tmpl")
	r.Static("/css", "server/resources/css")
	r.Static("/img", "server/resources/img")
	r.Static("/js", "server/resources/js")
    r.GET("/", func(c *gin.Context) {
        c.Header("Content-Type", "text/html")
        c.HTML(http.StatusOK, "index.tmpl", gin.H{})
    })

    r.POST("/search", func(c *gin.Context) {
        website := c.PostForm("website")
        c.String(200, "Website: %s", website)
    })

    return r
}