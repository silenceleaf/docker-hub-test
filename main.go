package main

import (
	//"flag"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes"
)

var (
	//kubeClient        *kubernetes.Clientset
	requestCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "test",
			Subsystem: "",
			Name:      "http_requests",
			Help:      "http response code counter",
		},
		[]string{"api", "status"},
	)
)

func main() {
	metricsRegistry := prometheus.NewRegistry()
	metricsRegistry.Register(requestCounterVec)
	// creates the in-cluster config
	// configFile := flag.String("kubeconfig", "c:\\dev\\go\\src\\org.junyan\\test\\kube-api-server.conf", "(optional) absolute path to the kubeconfig file")
	// config, err := clientcmd.BuildConfigFromFlags("", *configFile)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// creates the clientset
	// clientSet, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// kubeClient = clientSet

	router := gin.Default()
	router.GET("/get", get200)
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{ErrorHandling: promhttp.HTTPErrorOnError})))
	//router.GET("/test", testKubeSecretes)
	router.Run(":8888")
}

func hashAndTruncateLongName(input string, length int) string {
	md5Data := md5.Sum([]byte(input))
	base64Str := base64.StdEncoding.EncodeToString(md5Data[:md5.Size])
	fmt.Printf("hash name: %s", base64Str)
	return base64Str[:length]
}

func get200(c *gin.Context) {
	defer c.Next()
	requestCounterVec.WithLabelValues("test_api", "2xx").Inc()
	json := gin.H{
		"status":       "OK!",
		"hashLongName": hashAndTruncateLongName("this is a very long name which need to be hashed", 13),
	}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		json[pair[0]] = pair[1]
	}
	c.JSON(http.StatusCreated, json)
}

// func testKubeSecretes(c *gin.Context) {
// 	// pods, err := kubeClient.CoreV1().Pods("default").List(metav1.ListOptions{})
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }
// 	secrets, err := kubeClient.CoreV1().Secrets("default").Get("dbadmin.grafana", metav1.GetOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	c.IndentedJSON(http.StatusOK, gin.H{
// 		"endpoint": string(secrets.Data["endpoint"][:]),
// 		"port":     string(secrets.Data["port"][:]),
// 		"username": string(secrets.Data["username"][:]),
// 		"password": string(secrets.Data["password"][:]),
// 	})
// }
