package middleware

import (
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/seregant/golang_k8s_provisioning/config"
)

//type Default struct{}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func ValidateToken() gin.HandlerFunc {
	config := config.SetConfig()
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		if bearer != "" {
			// fmt.Println(bearer + "test")
			strSplit := strings.Split(bearer, " ")
			fmt.Println("token masuk ke validasi : ", strSplit[1])
			if strSplit[0] == "Bearer" && strSplit[1] != "" {
				secretKey := config.SecretKey
				token, err := jwt.Parse(strSplit[1], func(token *jwt.Token) (interface{}, error) {
					fmt.Println(token)
					if jwt.GetSigningMethod("HS512") != token.Method {
						return nil, fmt.Errorf("Unexpected signing method")
					}

					return []byte(secretKey), nil
				})
				fmt.Println("error : ", err)
				_, ok := token.Claims.(jwt.MapClaims)
				fmt.Println("validation result : ")
				fmt.Println(ok)
				fmt.Println(token.Valid)
				if ok && token.Valid && err == nil {
					c.Next()
				} else {
					c.JSON(401, gin.H{
						"status":  401,
						"message": "unauthorized",
					})
					c.AbortWithStatus(401)
				}
			} else {
				c.JSON(401, gin.H{
					"status":  401,
					"message": "unauthorized",
				})
				c.AbortWithStatus(401)
			}
		} else {
			c.JSON(401, gin.H{
				"status":  401,
				"message": "unauthorized",
			})
			c.AbortWithStatus(401)
		}
	}
}
