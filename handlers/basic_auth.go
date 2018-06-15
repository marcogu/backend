package handlers

import (
	"backend/models"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type VCodeForm struct {
	Mobile string `form:"mobile" binding:"numeric,len=11"`
}

type LoginForm struct {
	Mobile string `form:"mobile" binding:"numeric,len=11"`
	VCode  string `form:"vcode" binding:"numeric,len=4"`
}

func VCodeLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginForm LoginForm
		if err := c.ShouldBind(&loginForm); err == nil {
			if vcode, ok := utils.CACHE[loginForm.Mobile]; ok {
				if vcode == loginForm.VCode {
					c.JSON(http.StatusOK, gin.H{
						"message": "登录成功",
					})
				} else {
					c.JSON(http.StatusNotFound, gin.H{
						"message": "验证码错误",
					})
				}
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "验证码不存在，请重新发送",
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "无效手机号或验证码",
			})
		}
	}
}

func VCodeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vcodeForm VCodeForm
		if err := c.ShouldBind(&vcodeForm); err == nil {
			var user models.User
			db.Where(&models.User{Mobile: vcodeForm.Mobile}).First(&user)
			if user.ID == 0 {
				db.Create(&models.User{Mobile: vcodeForm.Mobile})
			}

			vcode := randomVerifyingCode()
			utils.CACHE[vcodeForm.Mobile] = vcode
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
