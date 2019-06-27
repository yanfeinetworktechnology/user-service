package controller

import (
	base_common "base_service/common"
	"net/http"
	"user_service/common"
	"user_service/model"

	"github.com/gin-gonic/gin"
)

var (
	Person = 1
)

// CertificationPerson 个人账号认证
// @Summary 个人账号认证
// @Description 个人账号认证
// @Tags 认证相关
// @Param token header string true "token"
// @Param user body model.PersonRequest true "个人认证信息"
// @Accept json
// @Produce json
// @Success 200 {object} model.Message
// @Router /certification/person [post]
func CertificationPerson(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	var personRequest model.PersonRequest
	if common.FuncHandler(c, c.BindJSON(&personRequest), nil, common.ParameterError) {
		return
	}

	db := base_common.GetMySQL()
	tx := db.Begin()

	var existUser model.User

	err := db.First(&existUser, userID).Error
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		tx.Rollback()
		return
	}

	// 之前没有实名信息，直接插入新的
	if existUser.InfoID == 0 {
		var newPerson model.Person
		newPerson.RealName = personRequest.RealName
		newPerson.Sex = personRequest.Sex
		newPerson.Phone = personRequest.Phone
		newPerson.Hometown = personRequest.Hometown

		err = db.Create(&newPerson).Error
		if common.FuncHandler(c, err, nil, common.DatabaseError) {
			tx.Rollback()
			return
		}

		updateData := map[string]interface{}{"role": Person, "info_id": newPerson.ID}
		err = tx.Model(&existUser).Updates(updateData).Error
		if common.FuncHandler(c, err, nil, common.DatabaseError) {
			tx.Rollback()
			return
		}
	} else {
		var existPerson model.Person
		err = db.First(&existPerson, existUser.InfoID).Error

		updateData := map[string]interface{}{"real_name": existPerson.RealName, "sex": existPerson.Sex, "homwtown": existPerson.Hometown, "phone": existPerson.Phone}
		err = tx.Model(&existPerson).Updates(updateData).Error
		if common.FuncHandler(c, err, nil, common.DatabaseError) {
			tx.Rollback()
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, model.Message{
		Msg: "信息完善成功",
	})
}
