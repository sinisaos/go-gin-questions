package controllers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sinisaos/questions/config"
	"github.com/sinisaos/questions/models"
	"golang.org/x/crypto/bcrypt"
)

var err error

func Login(c *gin.Context) {
	h := gin.H{}
	path := c.Query("next")
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	h["path"] = path
	session.Save()
	c.HTML(http.StatusOK, "login.tmpl.html", h)
}

func SignIn(c *gin.Context) {
	session := sessions.Default(c)
	host := c.Request.Header.Get("Host")
	next := c.PostForm("next")
	username := c.PostForm("username")
	password := c.PostForm("password")
	user := models.User{}
	config.DB.Where("username = ?", username).First(&user)

	if user.Username != username || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		session.AddFlash("Email or password incorrect")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return
	}

	config.DB.Model(&user).
		Where("users.username = ?", username).
		UpdateColumn("is_logged_in", gorm.Expr("is_logged_in + ?", 1))

	session.Set("user", username)
	session.Save()

	c.Redirect(http.StatusFound, host+next)
	return

}

func SignUp(c *gin.Context) {
	h := gin.H{}
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "signup.tmpl.html", h)
}

func SaveUser(c *gin.Context) {
	h := gin.H{}
	session := sessions.Default(c)
	user := models.User{}
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	passwordconfirm := c.PostForm("passwordconfirm")
	session.Set("username", username)
	session.Set("email", email)
	session.Set("pasword", password)
	h["username"] = username
	h["email"] = email
	h["password"] = password

	if username == "" {
		h["a"] = "Required field can't be empty!"
		c.HTML(http.StatusFound, "signup.tmpl.html", h)
		return
	}

	if email != "" {
		re := regexp.MustCompile(".+@.+\\..+")
		matched := re.Match([]byte(email))
		if matched == false {
			h["b"] = "Enter a valid email address!"
			c.HTML(http.StatusOK, "signup.tmpl.html", h)
			return
		}
	}

	if password != passwordconfirm {
		h["d"] = "The two password fields didn't match!"
		c.HTML(http.StatusOK, "signup.tmpl.html", h)
		return
	}

	var count int
	config.DB.Where("email = ?", email).Find(&user).Count(&count)
	if count != 0 {
		session.AddFlash("User already exists.")
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return
	}
	//create user
	user.Password, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Email = email
	user.Username = username

	if err := config.DB.Create(&user).Error; err != nil {
		session.AddFlash("User already exists.")
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return

	}
	//add first user as admin
	config.DB.Model(&user).
		Where("users.id = ?", 1).
		UpdateColumn("is_admin", gorm.Expr("is_admin + ?", 1))

	session.Set("user", username)
	session.Save()
	c.Redirect(http.StatusFound, "/")
	return
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	users := []models.User{}
	userId, _ := strconv.Atoi(id)
	config.DB.Delete(&users, userId)
	c.Redirect(http.StatusFound, "/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	usr := models.User{}
	user := session.Get("user")
	config.DB.Model(&usr).
		Where("users.username = ?", user).
		UpdateColumn("is_logged_in", gorm.Expr("is_logged_in - ?", 1))

	session.Delete("user")
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
	return
}

func Profile(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	c.Writer.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	c.Writer.Header().Set("Expires", "0")                                         // Proxies.
	session := sessions.Default(c)
	questions := []models.Question{}
	users := []models.User{}
	answers := []models.Answer{}
	id := c.Param("id")
	user := session.Get("user")
	var countQuestions int
	var countAnswers int
	var userId int

	config.DB.Find(&users)

	for _, v := range users {
		if v.Username == user {
			userId = v.Id
		}
	}

	config.DB.Where("users.id=?", id).Find(&users)
	config.DB.Preload("User").Where("questions.user_id=?", id).Find(&questions).Count(&countQuestions)
	config.DB.Preload("User").Where("answers.user_id=?", id).Find(&answers).Count(&countAnswers)
	session.Save()

	c.HTML(http.StatusOK, "profile.tmpl.html",
		gin.H{"questions": questions,
			"answers":        answers,
			"users":          users,
			"user":           user,
			"userId":         userId,
			"countQuestions": countQuestions,
			"countAnswers":   countAnswers})

}

func Admin(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	c.Writer.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	c.Writer.Header().Set("Expires", "0")                                         // Proxies.
	session := sessions.Default(c)
	questions := []models.Question{}
	users := []models.User{}
	answers := []models.Answer{}
	user := session.Get("user")
	var userId int
	var admin string

	config.DB.Find(&users)

	for _, v := range users {
		if v.Username == user {
			userId = v.Id
		}
		if v.IsAdmin == true {
			admin = v.Username
		}
	}

	if user != admin {
		c.Redirect(http.StatusSeeOther, "/")
	}

	config.DB.Find(&questions)
	config.DB.Find(&answers)
	session.Save()

	c.HTML(http.StatusOK, "admin.tmpl.html",
		gin.H{"questions": questions,
			"answers": answers,
			"users":   users,
			"userId":  userId,
			"admin":   admin,
			"user":    user,
		})
}

func RankUser(c *gin.Context) {
	type Result struct {
		Username string
		Total    int
	}
	results := []Result{}
	users := []models.User{}
	session := sessions.Default(c)
	user := session.Get("user")
	var userId int

	config.DB.Find(&users)

	for _, v := range users {
		if v.Username == user {
			userId = v.Id
		}
	}

	session.Save()
	config.DB.Raw("SELECT users.username, count(answers.is_accepted_answer) as total FROM users JOIN answers ON users.id = answers.user_id WHERE answers.is_accepted_answer = true GROUP BY users.id ORDER BY total DESC").Scan(&results)
	c.HTML(http.StatusOK, "ranking.tmpl.html",
		gin.H{"results": results,
			"user":   user,
			"userId": userId,
		})
}
