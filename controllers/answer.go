package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sinisaos/questions/config"
	"github.com/sinisaos/questions/models"
)

func SaveAnswer(c *gin.Context) {
	ip := c.Request.Header.Get("Referer")
	session := sessions.Default(c)
	u := c.PostForm("user")
	i := c.PostForm("id")
	body := c.PostForm("body")

	answerUserId, _ := strconv.Atoi(u)
	questionUserId, _ := strconv.Atoi(i)
	answers := models.Answer{
		UserID:     answerUserId,
		QuestionID: questionUserId,
		Body:       body,
	}

	config.DB.Save(&answers)
	config.DB.Exec("UPDATE questions SET answer_count = answer_count + 1 WHERE questions.id = ?", questionUserId)

	session.Save()
	c.Redirect(http.StatusFound, ip)
}

func AcceptAnswer(c *gin.Context) {
	answer := []models.Answer{}
	question := []models.Question{}
	id := c.PostForm("qid")
	ans := c.PostForm("aid")
	answerId, _ := strconv.Atoi(ans)
	questionId, _ := strconv.Atoi(id)
	config.DB.Model(&answer).
		Where("answers.id = ?", answerId).
		UpdateColumn("is_accepted_answer", gorm.Expr("is_accepted_answer + ?", 1))

	config.DB.Model(&question).
		Where("questions.id = ?", questionId).
		UpdateColumn("accepted_answer", gorm.Expr("accepted_answer + ?", 1))

	t := strconv.Itoa(questionId)
	c.Redirect(http.StatusFound, "/show/"+t)
}

func AnswerLikes(c *gin.Context) {
	answer := []models.Answer{}
	id := c.PostForm("qid")
	ans := c.PostForm("aid")
	answerId, _ := strconv.Atoi(ans)
	questionId, _ := strconv.Atoi(id)
	config.DB.Model(&answer).
		Where("answers.id = ?", answerId).
		UpdateColumn("likes", gorm.Expr("likes + ?", 1))

	t := strconv.Itoa(questionId)
	c.Redirect(http.StatusFound, "/show/"+t)
}

func AnswerDisLikes(c *gin.Context) {
	answer := []models.Answer{}
	id := c.PostForm("qid")
	ans := c.PostForm("aid")
	answerId, _ := strconv.Atoi(ans)
	questionId, _ := strconv.Atoi(id)
	config.DB.Model(&answer).
		Where("answers.id = ?", answerId).
		UpdateColumn("dis_likes", gorm.Expr("dis_likes + ?", 1))

	t := strconv.Itoa(questionId)
	c.Redirect(http.StatusFound, "/show/"+t)
}

/*
func AnswerDelete(c *gin.Context) {
	id := c.Params("id")
	ip := c.Request.Header.Get("Referer")
	answers := []models.Answer{}
	config.DB.Delete(&answers, id)
	c.Redirect(http.StatusFound, ip)
}*/
