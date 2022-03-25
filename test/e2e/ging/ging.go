package ging

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

type GinG struct {
	Describe      string
	KubeClientSet kubernetes.Interface
	KubeConfig    *rest.Config
}

func NewGinG(describe, kubeConfigPath string) *GinG {
	// generate instance
	g := &GinG{Describe: describe}
	// client-go
	conf, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		klog.Fatalf("[error when get kubeConfigPath]\n")
		panic(err.Error())

	}
	g.KubeConfig = conf

	kubeClient, err := kubernetes.NewForConfig(conf)
	if err != nil {
		klog.Fatalf("[error when new kubeClientSet]\n")
		panic(err.Error())
	}
	g.KubeClientSet = kubeClient
	return g
}

func (g *GinG) GetDescribe() string {
	return strings.Replace(CurrentGinkgoTestDescription().TestText, " ", "-", -1)
}
