// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Spiderpool

package framework

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"os"
	"time"
)

const (
	createTimeout = 3 * time.Minute
	deleteTimeout = 2 * createTimeout
)

type Framework struct {
	BaseName        string
	SystemNameSpace string
	KubeClientSet   *kubernetes.Clientset
	KubeConfig      *rest.Config
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
	f.KubeClientSet = kubeClient
	return f
}

func (f *Framework) CreatePod(podNamespace, podName, image string) (*corev1.Pod, error) {
	pod := podObject(podNamespace, podName, image)
	pod, err := f.KubeClientSet.CoreV1().Pods(pod.Namespace).Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	if err := WaitForPodReady(f.KubeClientSet, pod.Namespace, pod.Name, createTimeout); err != nil {
		return nil, err
	}

	pod, err = f.KubeClientSet.CoreV1().Pods(pod.Namespace).Get(context.Background(), pod.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (f *Framework) DeletePod(pod *corev1.Pod) error {
	if err := f.KubeClientSet.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	if err := WaitForPodToDisappear(f.KubeClientSet, pod.GetNamespace(), pod.GetName(), deleteTimeout); err != nil {
		return err
	}
	return nil
}

func podObject(podNamespace, podName, image string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: podNamespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            fmt.Sprintf("%s-container", podName),
					Command:         containerCmd(),
					Image:           image,
					ImagePullPolicy: "IfNotPresent",
				},
			},
		},
	}
}

func containerCmd() []string {
	return []string{"/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait"}
}
