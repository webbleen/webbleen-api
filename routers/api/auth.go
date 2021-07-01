package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/webbleen/go-gin/models"
	"github.com/webbleen/go-gin/pkg/e"
	"github.com/webbleen/go-gin/pkg/logging"
	"github.com/webbleen/go-gin/pkg/util"
)

type AuthLogin struct {
	Username string `json:"username" binding:"required" valid:"Required; MaxSize(50)"`
	Password string `json:"password" binding:"required" valid:"Required; MaxSize(50)"`
}

func Login(c *gin.Context) {
	var a AuthLogin
	c.BindJSON(&a)

	valid := validation.Validation{}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(a.Username, a.Password)
		if isExist {
			token, err := util.GenerateToken(a.Username, a.Password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token

				code = e.SUCCESS
			}

		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func Logout(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]interface{})

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
