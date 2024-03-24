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
        delimiter:="\\,\\"
        urls:= "a.com" + delimiter + "b.com"
        titles:="Site a" + delimiter + "Site b"
        data:= "aaaa" + delimiter + "bbbb"
        /*urls :=  make(map[int]string)
        urls[1] = "a.com"
        urls[2] = "b.com"
        titles :=  make(map[int]string)
        titles[1] = "Site a"
        titles[2] = "Site b"
        data :=  make(map[int]string)
        data[1] = "aaaa"
        data[2] = "bbbb"*/
        c.HTML(http.StatusOK, "search.tmpl", gin.H{
            "urls": urls,
            "titles": titles,
            "data": data,
            "website": website,
        })
    })

    return r
}