package main

import (
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"

    _ "github.com/mattn/go-sqlite3" // DB操作はGORMでやるので使わないという意味で - をつける
)

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/*.html")

	data := "This string is set in main.go"

    router.GET("/", func(ctx *gin.Context){
        ctx.HTML(200, "index.html", gin.H{"data": data})
    })

    router.Run()
}