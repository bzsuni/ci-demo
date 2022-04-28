// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Spider
package ip

import (
	"ci-demo/test/e2e/framework"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"net"
)

//func TestSetup(t *testing.T) {
//	RegisterFailHandler(Fail)
//	RunSpecs(t, "Setup Suite")
//}

//var _ = Describe("Kind Cluster Setup", func() {
//	f := framework.NewFramework("Setup")
//
//	Describe("Cluster Setup", func() {
//		It("List Spider Pods", func() {
//			_, err := f.KubeClientSet.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
//			Expect(err).NotTo(HaveOccurred())
//		})
//
//		It("Create Pod", func() {
//			pod := &corev1.Pod{
//				ObjectMeta: metav1.ObjectMeta{Name: "e2e-test"},
//				Spec: corev1.PodSpec{
//					Containers: []corev1.Container{
//						{
//							Name:    "samplepod",
//							Image:   "alpine",
//							Command: []string{"/bin/ash", "-c", "trap : TERM INT; sleep infinity & wait"},
//						},
//					},
//				},
//			}
//			_, err := f.KubeClientSet.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
//			Expect(err).NotTo(HaveOccurred())
//		})
//
//	})
//
//})

const (
	testNamespace = "default"
	testPodName   = "test-pod"
	testImg       = "alpine"
	testIPv4Range = "172.18.0.0/16" // from /tools/yamls/10-macvlan.conflist
)

var _ = Describe("Whereabouts functionality", func() {
	f := framework.NewFramework("test macvlan whereabouts")
	Context("Single pod tests", func() {
		var pod *corev1.Pod
		BeforeEach(func() {
			By("creating a pod with whereabouts net-attach-def")
			var err error

			// params podNamespace, podName, image string, label, annotations map[string]string
			pod, err = f.CreatePod(
				testNamespace,
				testPodName,
				testImg,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		//AfterEach(func() {
		//	By("deleting pod with whereabouts net-attach-def")
		//	Expect(f.DeletePod(pod)).To(Succeed())
		//})

		It("allocates a single pod within the correct IP range", func() {
			By("checking pod IP is within whereabouts IPAM range")
			podIP, err := getPodIP(pod)
			Expect(err).NotTo(HaveOccurred())
			Expect(inRange(testIPv4Range, podIP)).To(Succeed())
		})
	})
})

func getPodIP(pod *corev1.Pod) (string, error) {
	podIP := pod.Status.PodIP
	if len(podIP) == 0 {
		return "", fmt.Errorf("the pod does not have IP,maybe ip allocation err")
	}
	return podIP, nil
}

func inRange(cidr string, ip string) error {
	_, cidrRange, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	if cidrRange.Contains(net.ParseIP(ip)) {
		return nil
	}

	return fmt.Errorf("ip [%s] is NOT in range %s", ip, cidr)
}
