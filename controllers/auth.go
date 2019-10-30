package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/seregant/golang_k8s_provisioning/config"
	"github.com/seregant/golang_k8s_provisioning/database"
	"github.com/seregant/golang_k8s_provisioning/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthController is
type AuthController struct{}

//GenerateToken is
func (w *AuthController) GenerateToken(c *gin.Context) {
	username, _ := c.GetPostForm("user")
	password, _ := c.GetPostForm("password")
	fmt.Println(username)
	fmt.Println(password)
	emailPttrn := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	var data models.Pengguna
	var db = database.DbConnect()
	defer db.Close()

	if emailPttrn.MatchString(username) {
		db.Where("pengguna_email = ?", username).Where("pengguna_password = ?", password).First(&data)
	} else {
		db.Where("pengguna_username = ?", username).Where("pengguna_password = ?", password).First(&data)
	}

	if data.Username != "" {
		config := config.SetConfig()
		//generate new token
		secretKey := config.SecretKey
		sign := jwt.New(jwt.GetSigningMethod("HS512"))
		claims := sign.Claims.(jwt.MapClaims)
		expiredAt := time.Now().Add(time.Duration(config.TokenExpTime * 1000000000)).Unix()
		claims["iat"] = time.Now().Unix() //iat and exp is standart claims
		claims["exp"] = expiredAt
		claims["name"] = data.Nama
		claims["email"] = data.Email
		jwtToken, err := sign.SignedString([]byte(secretKey))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "success",
				"data": gin.H{
					"token":      jwtToken,
					"expired_at": expiredAt,
				},
			})
		}
	} else {
		c.JSON(401, gin.H{
			"status":  401,
			"message": "unauthorized username/password !",
		})
		c.AbortWithStatus(401)
	}
}
