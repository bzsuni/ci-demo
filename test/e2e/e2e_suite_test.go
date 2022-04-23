package e2e_test

import (
	"ci-demo/test/e2e/framework"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	// test cases
	_ "ci-demo/test/e2e/ip"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	f := framework.NewFramework("init")
	if f.BaseName == "init" {
		fmt.Println("ok")
	}
	return nil
}, func(data []byte) {})
