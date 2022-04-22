// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Spiderpool

package framework

import (
	"context"
	"fmt"
	nettypes "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"os"

	netclient "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned/typed/k8s.cni.cncf.io/v1"
)

type Framework struct {
	BaseName        string
	SystemNameSpace string
	KubeClientSet   kubernetes.Interface
	KubeConfig      *rest.Config
	NetClientSet    netclient.K8sCniCncfIoV1Interface
}

// NewFramework init Framework struct
func NewFramework(baseName string) *Framework {
	f := &Framework{BaseName: baseName}

	kubeconfigPath := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		klog.Fatal(err)
	}
	f.KubeConfig = cfg

	cfg.QPS = 1000
	cfg.Burst = 2000
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatal(err)
	}
	netClient, err := netclient.NewForConfig(cfg)
	if err != nil {
		klog.Fatal(err)
	}
	f.KubeClientSet = kubeClient
	f.NetClientSet = netClient

	return f
}

func (f *Framework) addNetAttachDef(netattach *nettypes.NetworkAttachmentDefinition) (*nettypes.NetworkAttachmentDefinition, error) {
	return f.NetClientSet.NetworkAttachmentDefinitions(netattach.ObjectMeta.Namespace).Create(context.TODO(), netattach, metav1.CreateOptions{})
}

func (f *Framework) delNetAttachDef(netattach *nettypes.NetworkAttachmentDefinition) error {
	return f.NetClientSet.NetworkAttachmentDefinitions(netattach.ObjectMeta.Namespace).Delete(context.TODO(), netattach.Name, metav1.DeleteOptions{})
}
