package cookie

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cookie manager suite")
}

var _ = Describe("Shopping cart", func() {
	Context("initially", func() {
		It("dummy test", func() {
			Expect(0).Should(BeZero())
		})
	})
})
