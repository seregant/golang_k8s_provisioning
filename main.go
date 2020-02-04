package main

import (
	"github.com/gin-gonic/gin"
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
	clusterController := new(controllers.Nodes)
	authController := new(controllers.AuthController)

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	api := router.Group("/api")
	{
		login := api.Group("login")
		{
			login.POST("/", authController.GenerateToken)
		}
		pengguna := api.Group("pengguna")
		pengguna.Use(middleware.CORSMiddleware())
		{
			pengguna.GET("/", penggunaController.GetAll)
			pengguna.POST("/add", penggunaController.Add)
			pengguna.GET("/u", middleware.ValidateToken(), penggunaController.GetDataPengguna)
		}

		cluster := api.Group("/clusters")
		{
			cluster.GET("/nodes", middleware.ValidateToken(), clusterController.GetNodesData)
		}
	}
	// var emailNotif []string
	// emailNotif = append(emailNotif, "indradota17@gmail.com")
	// message := "Halo, untuk mengakses Owncloud anda silahkan login ke url <br> " + conf.Domain + "/indraag/login"
	// controllers.SendNotif(emailNotif, message)
	router.Run(":" + conf.HttpPort)
}
