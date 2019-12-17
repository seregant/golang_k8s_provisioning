package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/seregant/golang_k8s_provisioning/config"
	"github.com/seregant/golang_k8s_provisioning/models"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Nodes struct{}

func (w *Nodes) GetNodesData(c *gin.Context) {
	cfg := config.SetConfig()
	bearer := c.Request.Header.Get("Authorization")
	strSplit := strings.Split(bearer, " ")

	secretKey := cfg.SecretKey
	token, _ := jwt.Parse(strSplit[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS512") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method")
		}

		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var admin float64
		admin = 1
		if claims["admin"] == admin {
			var response []models.NodesItem
			var parsed models.NodesData
			var nodeCapacity struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			}

			clientset := config.SetK8sClient()
			data, err := clientset.RESTClient().Get().AbsPath("/apis/metrics.k8s.io/v1beta1/nodes").DoRaw()
			if err != nil {
				log.Fatal(err)
			}

			json.Unmarshal(data, &parsed)
			nodeList, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
			if err == nil {
				if len(nodeList.Items) > 0 {
					nodeCount := 0
					for _, node := range parsed.Items {

						nodeList := &nodeList.Items[nodeCount]

						memQuantity := nodeList.Status.Allocatable[v1.ResourceMemory]
						cpuQuantity := nodeList.Status.Allocatable[v1.ResourceCPU]
						nodeCapacity.Memory = strconv.Itoa(int(memQuantity.Value())/1000) + "Ki"
						nodeCapacity.CPU = strconv.Itoa(int(cpuQuantity.Value()))

						response = append(response, models.NodesItem{
							Metadata:  node.Metadata,
							Timestamp: node.Timestamp,
							Window:    node.Window,
							Capacity:  nodeCapacity,
							Usage:     node.Usage,
							Condition: nodeList.Status.Conditions,
						})
					}
				} else {
					c.JSON(500, gin.H{
						"status":  500,
						"message": "Unable to read node list",
					})
					c.AbortWithStatus(500)
					fmt.Println("Unable to read node list")
				}
			} else {
				c.JSON(500, gin.H{
					"status":  500,
					"message": "Error while reading node list data",
				})
				c.AbortWithStatus(500)
				fmt.Println("Error while reading node list data: %v", err)
			}
			c.JSON(200, gin.H{
				"status":  "200",
				"message": "success",
				"data":    response,
			})
		} else {
			c.JSON(401, gin.H{
				"status":  401,
				"message": "unauthorized",
			})
			c.AbortWithStatus(401)
		}
	}
}
