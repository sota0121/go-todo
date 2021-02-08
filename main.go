package main

import (
	"strconv"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3" // DB操作はGORMでやるので使わないという意味で - をつける
)

// ========================
// References
// ========================
// https://qiita.com/hyo_07/items/59c093dda143325b1859
// https://qiita.com/tfrcm/items/e2a3d7ce7ab8868e37f7


// ========================
// Database
// ========================
// model schema
type Todo struct {
	gorm.Model // gormのモデルオブジェクトをベースに
	Text	string
	Status	string
}

func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("can not open database (dbInit) ")
	}
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

func dbInsert(text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("can not open database (dbInsert) ")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("can not open database (dbUpdate) ")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("can not open database (dbDelete) ")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

func dbGetAll() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("can not open database (dbGetAll) ")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("can not open database (dbGetOne) ")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}

// ========================
// Main Proc
// ========================
func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/*.html")

	dbInit()

	// Index
    router.GET("/", func(ctx *gin.Context){
        todos := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
    })

	// Create with form
	router.POST("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbInsert(text, status)
		ctx.Redirect(302, "/")
	})

	// Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"todo": todo})
	})

	// Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbUpdate(id, text, status)
		ctx.Redirect(302, "/")
	})

	// Check item deleted
	router.GET("/delete_check/:id", func(ctx *gin.Context){
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"todos": todo})
	})

	// Delete
	router.POST("/delete/:id", func(ctx *gin.Context){
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		ctx.Redirect(302, "/")
	})

    router.Run()
}