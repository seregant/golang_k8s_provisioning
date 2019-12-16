package config

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	DbHost       string
	DbPort       string
	DbUser       string
	DbPass       string
	DbName       string
	HttpPort     string
	SrvKey       string
	ServerIp     string
	TokenExpTime int64
	SecretKey    string
	Domain       string
	AdminEmail   string
	Debug        bool
}

func SetConfig() Config {
	var config Config

	//set configuration here
	config.DbHost = "127.0.0.1"
	config.DbPort = "3306"
	config.DbUser = "root"
	config.DbName = "provisioning_owncloud"
	config.DbPass = ""
	config.HttpPort = "1235"
	config.SrvKey = "Aw4s_g4l4k"
	config.ServerIp = "192.168.1.1"
	config.TokenExpTime = 1800
	config.SecretKey = "KJKJIds6sh"
	config.Domain = "stidust-web.site"
	config.AdminEmail = "indrasullivan17@gmail.com"
	config.Debug = true
	return config
}

func SetK8sClient() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", "./skripsi-cluster-kubeconfig.yaml")
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return clientset
}
