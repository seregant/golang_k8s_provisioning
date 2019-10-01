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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
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

func CheckClusterAvail() bool {

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
	// fmt.Println("cpu usage ", usageTotalCpu, " mCores of ", capacityTotalCpu*1000, "mCores")
	// fmt.Println("memory usage : ", usageTotalMem, " Ki of", capacityTotalMem, "Ki")

	if capacityTotalCpu-usageTotalCpu < 500 {
		return false
	} else if capacityTotalMem-usageTotalMem < 100000 {
		return false
	} else {
		return true
	}
}

func TestCluster() {
	config, err := clientcmd.BuildConfigFromFlags("", "./cluster-conf")
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	api := clientset.CoreV1()
	var ns, label, field string

	ns = "development"
	listOptions := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}

	pvcs, err := api.PersistentVolumeClaims(ns).List(listOptions)
	if err != nil {
		log.Fatal(err)
	}

	template := "%-32s%-8s%-8s\n"
	fmt.Printf(template, "NAME", "STATUS", "CAPACITY")
	for _, pvc := range pvcs.Items {
		quant := pvc.Spec.Resources.Requests[apiv1.ResourceStorage]
		fmt.Printf(
			template,
			pvc.Name,
			string(pvc.Status.Phase),
			quant.String())
	}

	pods, err := api.Pods(ns).List(listOptions)
	fmt.Println()
	fmt.Printf(template, "NAME", "STATUS", "HOST")
	for _, pod := range pods.Items {
		fmt.Printf(
			template,
			pod.Name,
			pod.Status.Phase,
			pod.Status.HostIP,
		)
	}
	template2 := "%-32s%-15s%-15s\n"
	nodes, err := api.Nodes().List(listOptions)
	fmt.Println()
	fmt.Printf(template2, "NAME", "CPU", "MEMORY")
	for _, node := range nodes.Items {
		fmt.Printf(
			template2,
			node.Name,
			node.Status.Capacity.Cpu(),
			node.Status.Capacity.Memory(),
		)
	}
}

func DeployOwnCloud() {
	clientset := config.SetK8sClient()
	deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nama-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "test-deploy", //-->from user data
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "test-deploy",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "owncloud-pengguna", //-->from variable by user ID
							Image: "owncloud/server",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "OWNCLOUD_DOMAIN",
									Value: "domain-variable", //-->from variable
								},
								{
									Name:  "OWNCLOUD_DB_TYPE",
									Value: "mysql",
								},
								{
									Name:  "OWNCLOUD_DB_HOST",
									Value: "103.56.205.130", //-->from variable
								},
								{
									Name:  "OWNCLOUD_DB_NAME",
									Value: "testing_indra_oc", //-->from variable
								},
								{
									Name:  "OWNCLOUD_DB_USERNAME",
									Value: "indra", //-->from variable
								},
								{
									Name:  "OWNCLOUD_DB_PASSWORD",
									Value: "indra123", //-->from variable
								},
								{
									Name:  "OWNCLOUD_ADMIN_USERNAME",
									Value: "indra", //-->from variable
								},
								{
									Name:  "OWNCLOUD_ADMIN_PASSWORD",
									Value: "indra123", //-->from variable
								},
								{
									Name:  "OWNCLOUD_REDIS_ENABLED",
									Value: "false", //-->from variable
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
		log.Fatal(err)
	}
	createService(deploy)
	fmt.Printf("Created deployment %q.\n", deploymentRes.GetObjectMeta().GetName())
}

func createService(a *appsv1.Deployment) {
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
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
			},
		},
	}

	service := clientset.CoreV1().Services("default")
	_, err := service.Create(serviceSpec)
	if err != nil {
		log.Fatal(err)
	}
}

// func createVol() {
// 	clientset := config.SetK8sClient()
// 	volSpec := &apiv1.PersistentVolume{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       "PersistentVolume",
// 			APIVersion: "v1",
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "pv-name",
// 		},
// 		Spec: apiv1.PersistentVolumeSpec{
// 			Capacity: apiv1.ResourceList{
// 				"storage": resource.MustParse("10Gi"),
// 			},
// 			AccessModes: []apiv1.PersistentVolumeAccessMode{
// 				"ReadWriteMany",
// 			},
// 			PersistentVolumeReclaimPolicy: "Retain",
// 		},
// 	}
// }

func int32Ptr(i int32) *int32 { return &i }
