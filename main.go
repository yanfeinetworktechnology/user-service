package main

import (
	"io"
	"user_service/model"

	base_common "base_service/common"
	user_common "user_service/common"

	"user_service/controller"

	"user_service/middleware"

	_ "user_service/docs"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func migrate(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin auto_increment=1")
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Person{})
}

// init 在 main 之前执行
func init() {
	// init config
	user_common.DefaultConfig()
	user_common.SetConfig()
	user_common.WatchConfig()

	// init logger
	base_common.InitLogger()

	// init Database
	db := base_common.InitMySQL()
	// 禁止在表名后面加s
	db.SingularTable(true)
	migrate(db)
}

// @title YANFEI-USER API
// @version 0.0.1
func main() {
	// Before init router
	if viper.GetBool("basic.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		// Redirect log to file
		gin.DisableConsoleColor()
		logFile := base_common.GetLogFile()
		defer logFile.Close()
		gin.DefaultWriter = io.MultiWriter(logFile)
	}

	r := gin.Default()
	// middleware
	r.Use(middleware.ErrorHandling())
	r.Use(middleware.MaintenanceHandling())
	r.Use(middleware.TokenHandling())

	// swagger router
	if viper.GetBool("basic.debug") {
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 路由
	r.POST("/user/login", controller.Login)
	r.POST("/certification/person", controller.CertificationPerson)

	r.Run("0.0.0.0:" + viper.GetString("basic.port"))
}
