package main

import (
	gin "github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/seregant/golang_k8s_provisioning/config"
	"github.com/seregant/golang_k8s_provisioning/controllers"
	"github.com/seregant/golang_k8s_provisioning/middleware"
)

var conf = config.SetConfig()

func main() {

	// database.DbInit()
	penggunaController := new(controllers.Pengguna)
	clusterController := new(controllers.Nodes)
	authController := new(controllers.AuthController)

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.POST("/login", authController.GenerateToken)
	api := router.Group("/api")
	{
		pengguna := api.Group("pengguna")
		{
			pengguna.GET("/", middleware.ValidateToken(), penggunaController.GetAll)
			pengguna.POST("/add", penggunaController.Add)
			pengguna.GET("/u", middleware.ValidateToken(), penggunaController.GetDataPengguna)
		}

		cluster := api.Group("/clusters")
		{
			cluster.GET("/nodes", clusterController.GetNodesData)
		}
	}
	router.Run(":" + conf.HttpPort)
}
