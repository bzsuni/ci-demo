// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Spiderpool
package ip

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/yaml"

	"ci-demo/test/e2e/framework"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const SpiderLabelSelector = "app.kubernetes.io/name: spiderpool"

var podTestYaml = `apiVersion: v1
kind: Pod
metadata:
  name: samplepod
spec:
  containers:
  - name: samplepod
    command: ["/bin/ash", "-c", "trap : TERM INT; sleep infinity & wait"]
    image: alpine
`

var _ = Describe("Kind Cluster Setup", func() {
	f := framework.NewFramework("Setup")

	Describe("Cluster Setup", func() {
		It("List Spider Pods", func() {
			pods, err := f.KubeClientSet.CoreV1().Pods("").List(context.Background(), v1.ListOptions{LabelSelector: SpiderLabelSelector})
			Expect(err).NotTo(HaveOccurred())

			for _, pod := range pods.Items {
				Expect(pod.Status.Phase).To(Equal(corev1.PodRunning))
			}
		})

		It("Create Pod", func() {
			pod := &corev1.Pod{}
			err := yaml.Unmarshal([]byte(podTestYaml), pod)
			Expect(err).NotTo(HaveOccurred())
			_, err = f.KubeClientSet.CoreV1().Pods("").Create(context.Background(), pod, v1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
