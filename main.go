package main

import (
	gin "github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/seregant/golang_k8s_provisioning/config"
	"github.com/seregant/golang_k8s_provisioning/controllers"
	"github.com/seregant/golang_k8s_provisioning/database"
	"github.com/seregant/golang_k8s_provisioning/middleware"
)

var conf = config.SetConfig()

func main() {

	database.DbInit()
	penggunaController := new(controllers.Pengguna)

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	api := router.Group("/api")
	{
		pengguna := api.Group("pengguna")
		{
			pengguna.GET("/", penggunaController.GetAll)
			pengguna.POST("/add", penggunaController.Add)
		}
	}
	router.Run(":" + conf.HttpPort)
}
