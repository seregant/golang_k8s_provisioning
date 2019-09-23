package controllers

import (
	"strconv"

	gin "github.com/gin-gonic/gin"
	"github.com/seregant/golang_k8s_provisioning/database"
	"github.com/seregant/golang_k8s_provisioning/models"
	"golang.org/x/crypto/bcrypt"
)

type Pengguna struct{}

func (w *Pengguna) GetAll(c *gin.Context) {
	var arr_response []models.PenggunaRes
	var data []models.Pengguna

	var db = database.DbConnect()
	defer db.Close()

	db.Find(&data)

	for _, data := range data {
		arr_response = append(arr_response, models.PenggunaRes{
			IDPengguna:  data.IDPengguna,
			Nama:        data.Nama,
			Alamat:      data.Alamat,
			Email:       data.Email,
			Username:    data.Username,
			Password:    data.Password,
			DBname:      data.DBname,
			DBuser:      data.DBuser,
			DBpass:      data.DBpass,
			ConfPath:    data.ClusterConf,
			StorageSize: data.StorageSize,
		})
	}
	c.JSON(200, gin.H{
		"status":  "200",
		"message": "success",
		"data":    arr_response,
	})
}

func (w *Pengguna) Add(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.JSON(405, gin.H{
			"status":  "405",
			"message": "method not allowed",
		})
	}
	var db = database.DbConnect()
	defer db.Close()
	var formData models.Pengguna
	formData.Nama, _ = c.GetPostForm("nama")
	formData.Alamat, _ = c.GetPostForm("alamat")
	formData.Email, _ = c.GetPostForm("email")
	formData.Username, _ = c.GetPostForm("username")
	formData.Password, _ = c.GetPostForm("password")
	storageSize, _ := c.GetPostForm("storage")
	formData.StorageSize, _ = strconv.Atoi(storageSize)

	formData.DBname = "db_" + formData.Username
	formData.DBuser = formData.Username
	dbPass, _ := bcrypt.GenerateFromPassword([]byte(formData.Username), 12)
	formData.DBpass = string(dbPass)

	db.Create(&formData)

	c.JSON(200, gin.H{
		"status":  "200",
		"message": "success",
	})
}
