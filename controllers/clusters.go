package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	gin "github.com/gin-gonic/gin"
	"github.com/seregant/golang_k8s_provisioning/config"
	"github.com/seregant/golang_k8s_provisioning/models"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

type Cluster struct{}

func (w *Cluster) GetNodesData(c *gin.Context) {
	config, err := clientcmd.BuildConfigFromFlags("", "./cluster-conf")
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").DoRaw()
	var parsed models.NodesData
	json.Unmarshal(data, &parsed)
	c.JSON(200, gin.H{
		"status":  "200",
		"message": "success",
		"data":    parsed,
	})
	CheckClusterAvail()
}

func Provisioning(dataUser models.Pengguna) bool {
	status := CheckClusterAvail()
	if status {
		if DeployDatabase(dataUser.DBpass, dataUser.DBname, dataUser.DBuser, dataUser.Username) {
			return DeployOwnCloud(dataUser.DBpass, dataUser.DBname, dataUser.DBuser, dataUser.Password, dataUser.Username, dataUser.Username+".domain.com")
		}
	}
	return false
}

func CheckClusterAvail() bool {
	fmt.Println("cek ketersediaan cluster")
	clientset := config.SetK8sClient()

	api := clientset.CoreV1()
	var label, field string

	listOptions := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}

	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").DoRaw()
	if err != nil {
		log.Fatal(err)
	}

	usageTotalCpu := 0
	usageTotalMem := 0
	var parsed models.NodesData
	json.Unmarshal(data, &parsed)
	for _, dataLoad := range parsed.Items {
		intCpu, _ := strconv.Atoi(strings.TrimSuffix(dataLoad.Usage.CPU, "m"))
		intMem, _ := strconv.Atoi(strings.TrimSuffix(dataLoad.Usage.Memory, "Ki"))
		usageTotalCpu = usageTotalCpu + intCpu
		usageTotalMem = usageTotalMem + intMem
		fmt.Println(intCpu)
	}

	capacityTotalCpu := 0
	capacityTotalMem := 0
	nodes, err := api.Nodes().List(listOptions)
	for _, node := range nodes.Items {
		intCpu, _ := strconv.Atoi(node.Status.Capacity.Cpu().String())
		intMem, _ := strconv.Atoi(strings.TrimSuffix(node.Status.Capacity.Memory().String(), "Ki"))
		capacityTotalCpu = capacityTotalCpu + intCpu
		capacityTotalMem = capacityTotalMem + intMem
	}
	//print usage to log
	// fmt.Println("cpu usage ", usageTotalCpu, " mCores of ", capacityTotalCpu*1000, "mCores")
	// fmt.Println("memory usage : ", usageTotalMem, " Ki of", capacityTotalMem, "Ki")

	if (capacityTotalCpu*1000)-usageTotalCpu < 500 {
		return false
	} else if capacityTotalMem-usageTotalMem < 100000 {
		return false
	} else {
		return true
	}
}

//jangan lupa buat log untuk setiap deployment yang telah dilakukan di db
func DeployOwnCloud(dbpass, dbname, dbuser, ocpass, ocuser, ocdomain string) bool {
	clientset := config.SetK8sClient()
	deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "owncloud-" + ocuser,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "oc-app-" + ocuser, //-->from user data
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "oc-app-" + ocuser,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "oc-usr" + ocuser, //-->from variable by user ID
							Image: "owncloud/server:10.2",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "owncloud",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8080,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "OWNCLOUD_DOMAIN",
									Value: ocdomain, //-->from variable
								},
								{
									Name:  "OWNCLOUD_DB_TYPE",
									Value: "mysql",
								},
								{
									Name:  "OWNCLOUD_DB_HOST",
									Value: "mysql-" + ocuser, //-->from variable
								},
								{
									Name:  "OWNCLOUD_DB_NAME",
									Value: dbname, //-->from variable
								},
								{
									Name:  "OWNCLOUD_DB_USERNAME",
									Value: dbuser, //-->from variable
								},
								{
									Name:  "OWNCLOUD_DB_PASSWORD",
									Value: dbpass, //-->from variable
								},
								{
									Name:  "OWNCLOUD_ADMIN_USERNAME",
									Value: ocuser, //-->from variable
								},
								{
									Name:  "OWNCLOUD_ADMIN_PASSWORD",
									Value: ocpass, //-->from variable
								},
								{
									Name:  "OWNCLOUD_REDIS_ENABLED",
									Value: "false", //-->from variable
								},
								// {
								// 	Name:  "HTTP_PORT",
								// 	Value: "80",
								// },
								// {
								// 	Name:  "HTTPS_PORT",
								// 	Value: "443",
								// },
							},
						},
					},
				},
			},
		},
	}
	fmt.Println("Creating deployment...")
	deploymentRes, err := deploymentClient.Create(deploy)
	if err != nil {
		return false
		log.Fatal(err)
	}
	fmt.Printf("Created deployment %q.\n", deploymentRes.GetObjectMeta().GetName())
	return createService(deploy, 8080)
}

func DeployDatabase(dbpass, dbname, dbuser, ocuser string) bool {
	clientset := config.SetK8sClient()
	deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mysql-" + ocuser,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mysql-app-" + ocuser, //-->from user data
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mysql-app-" + ocuser,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "mysql-usr" + ocuser, //-->from variable by user ID
							Image: "mysql:5.7",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "mysql",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 3306,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: dbpass, //-->from variable
								},
								{
									Name:  "MYSQL_DATABASE",
									Value: dbname, //-->from variable
								},
								{
									Name:  "MYSQL_USER",
									Value: dbuser, //-->from variable
								},
								{
									Name:  "MYSQL_PASSWORD",
									Value: dbpass, //-->from variable
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println("Creating deployment...")
	deploymentRes, err := deploymentClient.Create(deploy)
	if err != nil {
		return false
		log.Fatal(err)
	}
	fmt.Printf("Created deployment %q.\n", deploymentRes.GetObjectMeta().GetName())
	return createService(deploy, 3306)
}

func createService(a *appsv1.Deployment, port int32) bool {
	clientset := config.SetK8sClient()
	serviceSpec := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: a.ObjectMeta.Name, //-->from variable
		},
		Spec: apiv1.ServiceSpec{
			Type:     apiv1.ServiceTypeNodePort,
			Selector: a.Spec.Selector.MatchLabels,
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Protocol: apiv1.ProtocolTCP,
					Port:     port,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: port,
					},
				},
			},
		},
	}

	service := clientset.CoreV1().Services("default")
	_, err := service.Create(serviceSpec)
	if err != nil {
		return false
		log.Fatal(err)
	}
	return true
}

func CreateVol(pvName string, pvSize string) {
	var pvConf = `
{
	"apiVersion": "v1",
	"kind": "PersistentVolume",
	"metadata": {
		"name": "` + pvName + `"
	},
	"spec": {
		"capacity": {
			"storage": "` + pvSize + `Gi"
		},
		"accessModes": [{
			"ReadwriteMany"
		}]
		"persistentVolumeReclaimPolicy": "Retain",
		"nfs": {
			"server": "` + config.SetConfig().NfsServerIp + `",
			"path": "/opt/oc-data/users/` + pvName + `"
		},
		"claimRef": {
			"namespace": "development",
			"name": "` + pvName + `",
		}
	}
}
`
	var pvcConf = `
{
	"apiVersion": "v1",
	"kind": "PersistentVolumeClaim",
	"metadata": {
		"name": "` + pvName + `-pvc"
	},
	"spec": {
		"accessModes": [{
			"ReadWriteMany"
		}]
		"resources": {
			"resquests": {
				"storage": "` + pvSize + `Gi",
			}
		}
	}
}
`
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(pvConf), nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	deployment := obj.(*apiv1.PersistentVolume)

	fmt.Printf("%#v\n", deployment)
	obj2, _, err := decode([]byte(pvcConf), nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	deployment2 := obj2.(*apiv1.PersistentVolume)

	fmt.Printf("%#v\n", deployment2)
}

func VolumeTest() {
	var json = `
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-name
spec:
  capacity:
    storage: storage-size # This size is used to match a volume to a tenents claim
    accessModes:
      - ReadWriteMany # Access modes are defined below
    persistentVolumeReclaimPolicy: Retain # Reclaim policies are defined below
    nfs:
      server: 192.168.1.1
      path: nfs-server-path/user-folder
    claimRef:
      namespace: development
      name: pv-name
`
	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(json), nil, nil)
	if err != nil {
		fmt.Printf("error cuk", err)
	}

	// deployment := obj.(*v1beta1.Deployment)

	// fmt.Printf("%#v\n", obj.(type))
	switch o := obj.(type) {
	case *v1.PersistentVolume:
		fmt.Println("aproved")
	case *v1beta1.Role:
		// o is the actual role Object with all fields etc
	case *v1beta1.RoleBinding:
	case *v1beta1.ClusterRole:
	case *v1beta1.ClusterRoleBinding:
	case *v1.ServiceAccount:
	default:
		//o is unknown for us
		fmt.Println(o)
	}
}

func int32Ptr(i int32) *int32 { return &i }
