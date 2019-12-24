package controllers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	gin "github.com/gin-gonic/gin"
	"github.com/seregant/golang_k8s_provisioning/config"
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
			StorageSize: data.StorageSize,
			OcUrl:       data.OcUrl,
			IsAdmin:     data.IsAdmin,
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
	var penggunaQuery models.Pengguna
	var countUsername int
	var countEmail int
	var formData models.Pengguna
	formData.Nama, _ = c.GetPostForm("nama")
	formData.Alamat, _ = c.GetPostForm("alamat")
	formData.Email, _ = c.GetPostForm("email")
	formData.Username, _ = c.GetPostForm("username")
	formData.Password, _ = c.GetPostForm("password")
	storageSize, _ := c.GetPostForm("storage")
	formData.StorageSize, _ = strconv.Atoi(storageSize)

	if config.SetConfig().Debug {
		fmt.Println("DEBUG || TAMBAH USER : Req storage dari frontend : " + storageSize)
	}

	emailPttrn := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	db.Where("pengguna_username =  ?", formData.Username).Find(&penggunaQuery).Count(&countUsername)
	db.Where("pengguna_email = ?", formData.Email).Find(&penggunaQuery).Count(&countEmail)

	if emailPttrn.MatchString(formData.Email) {
		if countEmail == 0 {
			if countUsername == 0 {
				formData.DBname = "db_" + formData.Username
				formData.DBuser = formData.Username
				dbPass, _ := bcrypt.GenerateFromPassword([]byte(formData.Username), 12)
				formData.DBpass = string(dbPass)
				formData.OcUrl = formData.Username

				//jangan lupa notifikasi setelah provisioning berjalan
				if Provisioning(formData) {
					c.JSON(200, gin.H{
						"status":     "200",
						"message":    "success",
						"validation": "true",
						"details":    "Registrasi akun berhasil!",
					})
				} else {
					c.JSON(200, gin.H{
						"status":     "200",
						"message":    "success",
						"validation": "false",
						"details":    "Priovisioning gagal",
					})
				}
			} else {
				c.JSON(200, gin.H{
					"status":     "200",
					"message":    "success",
					"validation": "false",
					"details":    "Username sudah digunakan",
				})
			}
		} else {
			c.JSON(200, gin.H{
				"status":     "200",
				"message":    "success",
				"validation": "false",
				"details":    "Email sudah digunakan",
			})
		}
	} else {
		c.JSON(200, gin.H{
			"status":     "200",
			"message":    "success",
			"validation": "false",
			"details":    "Format email salah",
		})
	}
}

func (w *Pengguna) GetDataPengguna(c *gin.Context) {
	conf := config.SetConfig()
	var dataUser models.Pengguna
	var dataRes []models.PenggunaRes
	var db = database.DbConnect()
	defer db.Close()

	bearer := c.Request.Header.Get("Authorization")
	strSplit := strings.Split(bearer, " ")

	secretKey := config.SetConfig().SecretKey
	token, _ := jwt.Parse(strSplit[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS512") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method")
		}

		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if conf.Debug {
			fmt.Print("DEBUG || GET DATA USER : token yang diterima : ")
			fmt.Println(claims["admin"])
		}
		db.Where("pengguna_email = ?", claims["email"]).Find(&dataUser)
		dataRes = append(dataRes, models.PenggunaRes{
			IDPengguna:  dataUser.IDPengguna,
			Nama:        dataUser.Nama,
			Alamat:      dataUser.Alamat,
			Email:       dataUser.Email,
			Username:    dataUser.Username,
			Password:    dataUser.Password,
			DBname:      dataUser.DBname,
			DBuser:      dataUser.DBuser,
			DBpass:      dataUser.DBpass,
			StorageSize: dataUser.StorageSize,
			OcUrl:       conf.Domain + "/oc-client/" + dataUser.OcUrl,
			IsAdmin:     dataUser.IsAdmin,
		})
		c.JSON(200, gin.H{
			"status":  200,
			"message": "success",
			"data":    dataRes,
		})
	} else {
		fmt.Println("Invalid JWT Token")
	}
}
