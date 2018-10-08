package main

import (
	"flag"
	"net/http"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	kubeClient        *kubernetes.Clientset
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
	configFile := flag.String("kubeconfig", "c:\\dev\\go\\src\\org.junyan\\test\\kube-api-server.conf", "(optional) absolute path to the kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *configFile)
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	kubeClient = clientSet

	router := gin.Default()
	router.GET("/get", get200)
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{ErrorHandling: promhttp.HTTPErrorOnError})))
	router.GET("/test", testKubeSecretes)
	router.Run(":8080")
}

func get200(c *gin.Context) {
	defer c.Next()
	requestCounterVec.WithLabelValues("test_api", "2xx").Inc()
	c.JSON(http.StatusCreated, gin.H{"status": "OK!"})
}

func testKubeSecretes(c *gin.Context) {
	// pods, err := kubeClient.CoreV1().Pods("default").List(metav1.ListOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }
	secrets, err := kubeClient.CoreV1().Secrets("default").Get("dbadmin.grafana", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"endpoint": string(secrets.Data["endpoint"][:]),
		"port":     string(secrets.Data["port"][:]),
		"username": string(secrets.Data["username"][:]),
		"password": string(secrets.Data["password"][:]),
	})
}
