package config

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Host     string
	Port     string
	DbUser   string
	DbPass   string
	DbName   string
	HttpPort string
	SrvKey   string
}

func SetConfig() Config {
	var config Config

	//set configuration here
	config.Host = "cockroachdb"
	config.Port = "3306"
	config.DbUser = "root"
	config.DbName = "provisioning_owncloud"
	config.DbPass = ""
	config.HttpPort = "1235"
	config.SrvKey = "Aw4s_g4l4k"
	return config
}

func SetK8sClient() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", "./cluster-conf")
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return clientset
}
