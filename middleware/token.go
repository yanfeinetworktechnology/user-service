package middleware

import (
	base_common "base_service/common"
	"regexp"
	"user_service/common"
	"user_service/model"

	"github.com/gin-gonic/gin"
)

// TokenHandling 中间件，检查token
func TokenHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenWhiteList = []string{"/docs", "/user/login"}
		var certificationWhiteList = []string{"/certification/person"}
		var requestURL = c.Request.RequestURI

		for _, v := range tokenWhiteList {
			match, _ := regexp.MatchString(v, requestURL)
			if match {
				c.Next()
				return
			}
		}

		token := c.Request.Header.Get("token")

		if common.FuncHandler(c, token != "", true, common.NoToken) {
			c.Abort()
			return
		}

		j := common.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)

		if common.FuncHandler(c, err != common.ErrTokenExpired, true, common.TokenExpired) {
			c.Abort()
			return
		}

		if common.FuncHandler(c, err != common.ErrTokenInvalid, true, common.TokenInvalid) {
			c.Abort()
			return
		}

		// 检查是否需要认证，不需要则直接通过
		for _, v := range certificationWhiteList {
			match, _ := regexp.MatchString(v, requestURL)
			if match {
				c.Set("claims", claims)
				c.Next()
				return
			}
		}

		// 检查认证情况
		db := base_common.GetMySQL()
		userID := claims.UserID
		var existUser model.User
		err = db.First(&existUser, userID).Error

		if common.FuncHandler(c, err, nil, common.TokenInvalid) {
			c.Abort()
			return
		}
		if common.FuncHandler(c, existUser.InfoID != 0, true, common.NoCertification) {
			c.Abort()
			return
		}

		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
		c.Next()
	}
}
