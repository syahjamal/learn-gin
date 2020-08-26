package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Models
type Article struct {
	gorm.Model //Wajib untuk db apapun
	Title      string
	Slug       string `gorm:"unique_index"` //Kolom data ramah url dan urlnya cth: "judul-pertama"
	Desc       string `sql:"type:text;"`
}

//Variable global buat baca db
var DB *gorm.DB

func main() {
	var err error

	//koneksi db
	DB, err = gorm.Open("mysql", "root:@/go_gin_gorm?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer DB.Close()

	//Migrasi Model Article dari gorm
	DB.AutoMigrate(&Article{})

	router := gin.Default()

	//Grouping router agar rapih dan jika ada perubahan mudah untuk trace
	v1 := router.Group("/api/v1/")
	{
		article := v1.Group("/article")
		{
			article.GET("/", getHome)
			article.GET("/:slug", getArticle)
			article.POST("/", postArticle)
		}

		// users := v1.Group("/users")
		// {
		// 	users.GET("/", getUser)
		// }
	}
	router.Run()
}

func getHome(c *gin.Context) {

	items := []Article{}
	DB.Find(&items)

	c.JSON(200, gin.H{
		"status": "berhasil ke halaman home",
		"data":   items,
	})
}

func getArticle(c *gin.Context) {
	//parameter
	slug := c.Param("slug")

	var item Article

	//Query di gorm = Select * from table where slug = "slug"
	if DB.First(&item, "slug = ?", slug).RecordNotFound() {
		c.JSON(404, gin.H{"status": "error", "message": "record not found"})
		c.Abort() //Batalin request
		return
	}

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   item,
	})
}

func postArticle(c *gin.Context) {
	item := Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Slug:  slug.Make(c.PostForm("title")),
	}

	//Mencegah slug sama, maka generate random slug

	DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "berhasil post",
		"data":   item,
	})
}
