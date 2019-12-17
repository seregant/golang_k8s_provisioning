package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/seregant/golang_k8s_provisioning/config"
	"github.com/seregant/golang_k8s_provisioning/models"
)

var conf = config.SetConfig()

func DbConnect() *gorm.DB {
	var addr = conf.DbUser + ":" + conf.DbPass + "@tcp(" + conf.DbHost + ":" + conf.DbPort + ")/" + conf.DbName + "?charset=utf8&parseTime=True&loc=Local"
	if conf.Debug {
		fmt.Print("LOADED DB CONNECTION : ")
		fmt.Println(addr)
	}
	db, err := gorm.Open("mysql", addr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func DbInit() {

	var db = DbConnect()
	defer db.Close()

	db.Exec("CREATE DATABASE " + conf.DbName)
	if config.SetConfig().Debug {
		fmt.Println("Creating tables...")
	}

	db.AutoMigrate(&models.Pengguna{})
	db.AutoMigrate(&models.ClusterLoad{})
	db.AutoMigrate(&models.ClusterNode{})
}
