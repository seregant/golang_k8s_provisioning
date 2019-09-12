package controllers

import (
	"fmt"
	"github.com/seregant/golang_k8s_provisioning/database"
	"github.com/seregant/golang_k8s_provisioning/models"
	gin "github.com/gin-gonic/gin"
)

type Pengguna struct{}

func (w *Pengguna) GetAll(c *gin.Context){
	var arr_response []models.PenggunaRes
	var data []models.Pengguna

	var db = database.DbConnect()
	defer db.Close()

	db.Find(&data)

	for _, data := range data {
		arr_response = append(arr_response, models.PenggunaRes{
			IDPengguna: data.IDPengguna,
			Nama: data.Nama,
			Alamat: data.Alamat,
			Email: data.Email,
			Username: data.Usename,
			Password: data.Password,
			DBname: data.DBname,
			DBuser: data.DBuser,
			DBpass: data.DBpass.
			ConfPath: data.ClusterConf,
			StorageSize: data.StorageSize,
		 })
	}
}