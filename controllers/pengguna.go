package controllers

import (
	gin "github.com/gin-gonic/gin"
	"github.com/seregant/golang_k8s_provisioning/database"
	"github.com/seregant/golang_k8s_provisioning/models"
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
			Username:    data.Usename,
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
	var db = database.DbConnect()
	defer db.Close()

}
