package main

import (
	"fmt"
	"golang-training/week2/config"
	"golang-training/week2/provider"
	"net/http"
	"strconv"

	"github.com/microcosm-cc/bluemonday"

	"github.com/gin-gonic/gin"
)

// func main() {
// r.GET("/ping", func(c *gin.Context) {
// 	c.JSON(200, gin.H{
// 		"message": "pong",
// 	})
// })

// counter := 0
// r.GET("/increase", func(c *gin.Context) {
// 	counter++
// 	c.JSON(200, gin.H{
// 		"counter": counter,
// 	})
// })
// }

type AppError int

const (
	DB_ERROR    AppError = 1
	INPUT_ERROR AppError = 2
)

type Blog struct {
	ID      int
	Title   string `form:"title" binding:"min=10,max=100,required"`
	Content string `form:"content" binding:"min=15,max=1000"`
}

// TableName overrides the table name used by User to `profiles`
func (b Blog) TableName() string {
	return "blog"
}

func main() {
	cfg := config.NewConfig()
	rp := provider.MustBuildResourceProvider(cfg)
	db := rp.GetDB()

	r := gin.Default()

	r.GET("/health/live", func(c *gin.Context) {
		c.String(200, "ok")
	})

	blogRouter := r.Group("blog")
	// middleware
	blogRouter.Use(func(c *gin.Context) {
		lang, err := c.Cookie("lang")
		if err != nil || lang == "" {
			lang = "en"
		}
		c.Set("lang", lang)
	})
	blogRouter.Use(func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	})
	// handler
	// CREATE
	blogRouter.POST("/", func(c *gin.Context) {
		lang, _ := c.Get("lang")
		fmt.Println("lang", lang)
		c.Header("X-CUSTOM-LANG", lang.(string))

		blog := Blog{}
		if err := c.ShouldBindJSON(&blog); err != nil {
			// c.AbortWithStatus(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    INPUT_ERROR,
				"message": "Invalid format data",
				"detail":  err.Error(),
			})
			return
		}

		// sanitize
		p := bluemonday.UGCPolicy()
		blog.Title = p.Sanitize(blog.Title)
		blog.Content = p.Sanitize(blog.Content)

		if err := db.Create(&blog).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    DB_ERROR,
				"message": "Cannot create blog",
				"detail":  err.Error(),
			})
			return
		}
		c.JSON(200, blog)
	})

	// GET
	blogRouter.GET("/", func(c *gin.Context) {
		blog := Blog{}
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    INPUT_ERROR,
				"message": "Id is required",
			})
			return
		}
		if err := db.Where("id = ?", id).First(&blog).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    DB_ERROR,
				"message": "Cannot find blog",
				"detail":  err.Error(),
			})
			return
		}
		c.JSON(200, blog)
	})

	// LIST
	blogRouter.GET("/list", func(c *gin.Context) {
		blog := []Blog{}
		strpage := c.DefaultQuery("page", "1")
		strsize := c.DefaultQuery("size", "2")
		page, _ := strconv.ParseInt(strpage, 10, 32)
		size, _ := strconv.ParseInt(strsize, 10, 32)

		builder := db.Limit(int(size)).Offset(int(page)).Order("id DESC")

		q := c.DefaultQuery("q", "")
		if q != "" {
			// khong tot, nen tao 1 field search (fulltext index)
			// analyzer => filter nhugn ki tu khong can thiet (html, dup, ...)
			builder.Where("title LIKE '%?%' OR content LIKE '%?%'", q, q)
		}

		if err := builder.Find(&blog).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    DB_ERROR,
				"message": "Cannot find blog",
				"detail":  err.Error(),
			})
			return
		}
		c.JSON(200, blog)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
