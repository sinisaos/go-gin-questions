package main

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sinisaos/questions/config"
	"github.com/sinisaos/questions/controllers"
	"github.com/sinisaos/questions/models"
)

var err error

func main() {

	config.DB, err = gorm.Open("mysql", "root:1234@/gindb?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println("Status: ", err)
	}
	defer config.DB.Close()
	config.DB.LogMode(true)
	config.DB.AutoMigrate(&models.User{}, &models.Question{}, &models.Answer{}, &models.Tag{})
	config.DB.Model(&models.Question{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	config.DB.Model(&models.Answer{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "CASCADE")
	config.DB.Model(&models.Answer{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	r := gin.Default()

	r.Use(gin.Logger())
	r.LoadHTMLGlob("views/*.tmpl.html")

	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	store := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("qussession", store))

	//questions
	r.GET("/", controllers.AllQuestions)
	r.GET("/unsolved", controllers.UnsolvedQuestions)
	r.GET("/solved", controllers.SolvedQuestions)
	r.GET("/viewed", controllers.MostViewedQuestions)
	r.GET("/oldest", controllers.OldestQuestions)
	r.GET("/search", controllers.SearchQuestions)
	r.GET("/show/:id", controllers.ShowQuestion)
	r.GET("/create", controllers.CreateQuestion)
	r.GET("/edit/:id", controllers.EditQuestion)
	r.POST("/update/:id", controllers.UpdateQuestion)
	r.POST("/delete/:id", controllers.DeleteQuestion)
	r.POST("/savequestion", controllers.SaveQuestion)
	r.POST("/questionlikes", controllers.QuestionLikes)
	r.POST("/saveanswer", controllers.SaveAnswer)
	//answers
	r.POST("/acceptanswer", controllers.AcceptAnswer)
	r.POST("/answerlikes", controllers.AnswerLikes)
	r.POST("/answerdislikes", controllers.AnswerDisLikes)
	r.GET("/answeredit/:id", controllers.EditAnswer)
	r.POST("/answerupdate/:id", controllers.UpdateAnswer)
	r.POST("/answerdelete/:id", controllers.AnswerDelete)
	//tags
	r.GET("/tags/:id", controllers.Tags)
	r.GET("/tagedit/:id", controllers.EditTag)
	r.POST("/tagupdate", controllers.UpdateTag)
	r.GET("/categories", controllers.Categories)
	//users
	r.GET("/signup", controllers.SignUp)
	r.POST("/save", controllers.SaveUser)
	r.POST("/deleteuser/:id", controllers.DeleteUser)
	r.GET("/logout", controllers.Logout)
	r.GET("/login", controllers.Login)
	r.POST("/signin", controllers.SignIn)
	r.GET("/profile/:id", controllers.Profile)
	r.GET("/admin", controllers.Admin)
	r.GET("/rank", controllers.RankUser)

	r.Run(":8080")
}
