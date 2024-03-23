package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func Router() *gin.Engine {
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.Header("Content-Type", "text/html")
        c.HTML(http.StatusOK, "index.html", gin.H{})
    })

    return r
}