package controller

import (
	base_common "base_service/common"
	"net/http"
	"user_service/common"
	"user_service/model"

	"github.com/gin-gonic/gin"
)

// Login 统一用户登录
// @Summary 统一用户登录
// @Description 统一用户登录
// @Tags 用户相关
// @Param user body model.LoginRequest true "用户登录信息"
// @Accept json
// @Produce json
// @Success 200 {object} model.Message
// @Router /user/login [post]
func Login(c *gin.Context) {
	var loginRequest model.LoginRequest
	if common.FuncHandler(c, c.BindJSON(&loginRequest), nil, common.ParameterError) {
		return
	}
	db := base_common.GetMySQL()
	var existUser model.User
	var token string

	err := db.Where("user_name = ? AND password = ?", loginRequest.UserName, loginRequest.Password).First(&existUser).Error

	if common.FuncHandler(c, err, nil, common.LoginError) {
		return
	}
	// 成功
	token, err = common.CreateToken(existUser.ID)
	if common.FuncHandler(c, err, nil, common.SystemError) {
		return
	}

	c.JSON(http.StatusOK, model.Message{
		Data: token,
		Msg:  "登录成功",
	})
}
