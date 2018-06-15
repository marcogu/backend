package handlers

import (
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type LoginForm struct {
	Mobile string `form:"mobile" binding:"numeric,len=11"`
}

func VCodeLoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func VCodeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginForm LoginForm
		if err := c.ShouldBind(&loginForm); err == nil {
			vcode := randomVerifyingCode()
			utils.CACHE[loginForm.Mobile] = vcode
			c.JSON(http.StatusOK, gin.H{
				"vcode": vcode,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "无效手机号",
			})
		}
	}
}

func PasswordLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ModifyLoginPasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func randomVerifyingCode() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%04s", strconv.Itoa(rand.Intn(10000)))
}
