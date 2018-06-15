package handlers

import (
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func VCodeLoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func VCodeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		mobile := c.PostForm("mobile")
		if len(strings.TrimSpace(mobile)) == 11 {
			vcode := randomVerifyingCode()
			utils.CACHE[mobile] = vcode
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
