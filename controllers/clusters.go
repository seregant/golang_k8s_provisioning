package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/seregant/golang_k8s_provisioning/config"
	"github.com/seregant/golang_k8s_provisioning/database"
	"github.com/seregant/golang_k8s_provisioning/models"
	"k8s.io/apimachinery/pkg/api/resource"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1b1ex "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

func Provisioning(dataUser models.Pengguna) bool {
	// pvName := "volume-" + dataUser.Username
	// dbPvName := "mysql-pv-" + dataUser.Username
	var db = database.DbConnect()
	defer db.Close()
	var conf = config.SetConfig()
	status := CheckClusterAvail()
	storageSize := strconv.Itoa(dataUser.StorageSize)
	if status {
		if /*CreateVol(dbPvName, strconv.Itoa(dataUser.StorageSize))*/ true {
			if DeployDatabase(
				dataUser.DBpass,
				dataUser.DBname,
				dataUser.DBuser,
				dataUser.Username,
			) {
				if /*CreateVol(pvName, strconv.Itoa(dataUser.StorageSize))*/ true {
					if DeployOwnCloud(
						dataUser.DBpass,
						dataUser.DBname,
						dataUser.DBuser,
						dataUser.Password,
						dataUser.Username,
						config.SetConfig().Domain+"/oc-client/"+dataUser.Username,
						storageSize,
					) {
						db.Create(&dataUser)
						var emailNotif []string
						emailNotif = append(emailNotif, dataUser.Email)
						message := "Halo, untuk mengakses Owncloud anda silahkan login ke url " + conf.Domain + "/login"

						if IngressApply() {
							return sendNotif(emailNotif, message)
						}
					}
				}
			}
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
func DeployOwnCloud(dbpass, dbname, dbuser, ocpass, ocuser, ocdomain, ocstorage string) bool {
	fmt.Println("resquest storage : " + ocstorage)
	pvName := "volume-" + ocuser
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
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      pvName,
									MountPath: "/var/www/owncloud/data",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: pvName,
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-" + pvName,
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
	return createService(deploy, 8080)
}

func DeployDatabase(dbpass, dbname, dbuser, ocuser string) bool {
	pvName := "mysql-pv-" + ocuser
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
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      pvName,
									MountPath: "/var/lib/mysql",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: pvName,
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-" + pvName,
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

func CreateVol(pvName string, pvSize string) bool {
	cmd := exec.Command("mkdir", "/opt/oc-data/users/"+pvName)
	if err := cmd.Run(); err != nil {
		fmt.Print("Creating nfs folder error : ")
		fmt.Println(err)
		return false
	}
	clientset := config.SetK8sClient()
	k8sApi := clientset.CoreV1()
	var volAccModes []apiv1.PersistentVolumeAccessMode
	volAccModes = append(volAccModes, "ReadWriteMany")
	volSpec := &apiv1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: pvName,
		},
		Spec: apiv1.PersistentVolumeSpec{
			Capacity: apiv1.ResourceList{
				"storage": resource.MustParse(pvSize + "Gi"),
			},
			AccessModes:                   volAccModes,
			PersistentVolumeReclaimPolicy: apiv1.PersistentVolumeReclaimRetain,
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				NFS: &apiv1.NFSVolumeSource{
					Server: config.SetConfig().ServerIp,
					Path:   "/opt/oc-data/users/" + pvName,
				},
			},
			ClaimRef: &apiv1.ObjectReference{
				Name:      "pvc-" + pvName,
				Namespace: "default",
			},
		},
	}
	pvcSpec := &apiv1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "pvc-" + pvName,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: volAccModes,
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					"storage": resource.MustParse(pvSize + "Gi"),
				},
			},
		},
	}
	fmt.Println("Creating volume and it's pvc....")
	_, err := k8sApi.PersistentVolumes().Create(volSpec)
	if err != nil {
		fmt.Print("Creating volume error : ")
		fmt.Println(err)
		return false
	} else {
		fmt.Println("add volume " + pvName + " succeed")
		_, err2 := k8sApi.PersistentVolumeClaims(apiv1.NamespaceDefault).Create(pvcSpec)
		if err2 != nil {
			fmt.Print("Creating volume claim error : ")
			fmt.Println(err2)
			return false
		} else {
			fmt.Println("add volume claim pvc-" + pvName + " succeed")
			return true
		}
	}
}

func IngressApply() bool {
	fmt.Println("Updating ingress configuration..")
	var dataUser []models.Pengguna
	var db = database.DbConnect()
	defer db.Close()

	clientset := config.SetK8sClient()
	ingressClient := clientset.ExtensionsV1beta1().Ingresses(apiv1.NamespaceDefault)
	var ingressRules []v1b1ex.IngressRule
	var routeList []v1b1ex.HTTPIngressPath

	db.Find(&dataUser)

	for _, data := range dataUser {
		routeList = append(routeList, v1b1ex.HTTPIngressPath{
			Backend: v1b1ex.IngressBackend{
				ServiceName: "owncloud-" + data.Username,
				ServicePort: intstr.FromInt(8080),
			},
			Path: "/oc-client/" + data.OcUrl + "/?(.*)",
		})
	}

	ingressRules = append(ingressRules, v1b1ex.IngressRule{
		Host: config.SetConfig().Domain,
		IngressRuleValue: v1b1ex.IngressRuleValue{
			HTTP: &v1b1ex.HTTPIngressRuleValue{
				Paths: routeList,
			},
		},
	})
	ingressSpec := &v1b1ex.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "oc-nginx-ingress",
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target":     "/$1",
				"ingress.kubernetes.io/ssl-redirect":             "false",
				"kubernetes.io/ingress.class":                    "nginx",
				"nginx.ingress.kubernetes.io/force-ssl-redirect": "false",
			},
		},
		Spec: v1b1ex.IngressSpec{
			Rules: ingressRules,
		},
	}

	_, err := ingressClient.Create(ingressSpec)
	if err != nil {
		fmt.Print("Create ingress : ")
		fmt.Println(err)
		_, errUpdate := ingressClient.Update(ingressSpec)
		if errUpdate != nil {
			fmt.Print("Update ingress : ")
			fmt.Println(err)
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}
func int32Ptr(i int32) *int32 { return &i }
