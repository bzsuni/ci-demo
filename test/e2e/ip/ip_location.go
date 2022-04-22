package ip

import (
	"encoding/json"
	"fmt"
	nettypes "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"net"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"ci-demo/test/e2e/framework"
)

const (
	testNetworkName = "net-macvlan"
	testNamespace   = "default"
	testImage       = "quay.io/dougbtv/alpin:latest"
	ipv4TestRange   = "192.210.0.0/16"
	testPodName     = "test-pod"
)

var _ = Describe("Whereabouts functionality", func() {
	Context("test macvlan whereabouts", func() {
		var (
			netAttachDef *nettypes.NetworkAttachmentDefinition
		)
		f := framework.NewFramework("test macvlan whereabouts")

		BeforeEach(func() {
			netAttachDef = macvlanNetworkWithWhereaboutsIPAMNetwork()

			By("creating a NetworkAttachmentDefinition for whereabouts")
			_, err := f.AddNetAttachDef(netAttachDef)
			Expect(err).NotTo(HaveOccurred())
		})

		//AfterEach(func() {
		//	Expect(f.DelNetAttachDef(netAttachDef)).To(Succeed())
		//})

		Context("Single pod tests", func() {
			var (
				pod *corev1.Pod
			)
			BeforeEach(func() {
				By("creating a pod with whereabouts net-attach-def")
				var err error
				pod, err = f.CreatePod(
					testNamespace,
					testPodName,
					testImage,
					podTierLabel(testPodName),
					podNetworkSelectionElements(testNetworkName),
				)
				Expect(err).NotTo(HaveOccurred())
			})

			//AfterEach(func() {
			//	By("deleting pod with whereabouts net-attach-def")
			//	Expect(f.DeletePod(pod)).To(Succeed())
			//})

			It("allocates a single pod within the correct IP range", func() {
				By("checking pod IP is within whereabouts IPAM range")
				secondaryIfaceIP, err := secondaryIfaceIPValue(pod)
				Expect(err).NotTo(HaveOccurred())
				Expect(inRange(ipv4TestRange, secondaryIfaceIP)).To(Succeed())
			})
		})
	})
})

func clusterConfig() (*rest.Config, error) {
	const kubeconfig = "KUBECONFIG"

	kubeconfigPath, found := os.LookupEnv(kubeconfig)
	if !found {
		return nil, fmt.Errorf("must provide the path to the kubeconfig via the `KUBECONFIG` env variable")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func podTierLabel(podTier string) map[string]string {
	const tier = "tier"
	return map[string]string{tier: podTier}
}

func podNetworkSelectionElements(networkNames ...string) map[string]string {
	return map[string]string{
		nettypes.NetworkAttachmentAnnot: strings.Join(networkNames, ","),
	}
}

// Returns a network attachment definition object configured by provided parameters
func generateNetAttachDefSpec(name, namespace, config string) *nettypes.NetworkAttachmentDefinition {
	return &nettypes.NetworkAttachmentDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "NetworkAttachmentDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: nettypes.NetworkAttachmentDefinitionSpec{
			Config: config,
		},
	}
}

func macvlanNetworkWithWhereaboutsIPAMNetwork() *nettypes.NetworkAttachmentDefinition {
	macvlanConfig := `{
        "cniVersion": "0.3.0",
      	"disableCheck": true,
        "plugins": [
            {
                "type": "macvlan",
              	"master": "eth0",
              	"mode": "bridge",
              	"ipam": {
                    "type": "whereabouts",
                    "leader_lease_duration": 1500,
                    "leader_renew_deadline": 1000,
                    "leader_retry_period": 500,
                    "range": "10.10.0.0/16",
                    "log_level": "debug",
                    "log_file": "/tmp/wb"
              	}
            }
        ]
    }`
	return generateNetAttachDefSpec(testNetworkName, testNamespace, macvlanConfig)
}

func containerCmd() []string {
	return []string{"/bin/ash", "-c", "trap : TERM INT; sleep infinity & wait"}
}

func filterNetworkStatus(
	networkStatuses []nettypes.NetworkStatus, predicate func(nettypes.NetworkStatus) bool) *nettypes.NetworkStatus {
	for i, networkStatus := range networkStatuses {
		if predicate(networkStatus) {
			return &networkStatuses[i]
		}
	}
	return nil
}

func secondaryIfaceIPValue(pod *corev1.Pod) (string, error) {
	podNetStatus, found := pod.Annotations[nettypes.NetworkStatusAnnot]
	if !found {
		return "", fmt.Errorf("the pod must feature the `networks-status` annotation")
	}

	var netStatus []nettypes.NetworkStatus
	if err := json.Unmarshal([]byte(podNetStatus), &netStatus); err != nil {
		return "", err
	}

	secondaryInterfaceNetworkStatus := filterNetworkStatus(netStatus, func(status nettypes.NetworkStatus) bool {
		return status.Interface == "net1"
	})

	if len(secondaryInterfaceNetworkStatus.IPs) == 0 {
		return "", fmt.Errorf("the pod does not have IPs for its secondary interfaces")
	}

	return secondaryInterfaceNetworkStatus.IPs[0], nil
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
