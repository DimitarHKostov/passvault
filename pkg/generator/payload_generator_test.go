package generator

import (
	"testing"
	"time"

	logMock "passvault/pkg/log/mocks"

	"github.com/golang/mock/gomock"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cookie manager suite")
}

var _ = Describe("Shopping cart", func() {
	var (
		ctrl             *gomock.Controller
		logManagerMock   *logMock.MockLogManagerInterface
		payloadGenerator PayloadGeneratorInterface
	)

	const (
		expirationTime = 10 * time.Minute
	)

	var _ = BeforeEach(func() {
		ctrl = gomock.NewController(ginkgo.GinkgoT())
		logManagerMock = logMock.NewMockLogManagerInterface(ctrl)

		payloadGenerator = &PayloadGenerator{logManager: logManagerMock}
	})

	var _ = AfterEach(func() {
		ctrl.Finish()
	})

	Context("GeneratePayload", func() {
		It("should succeed to create a new payload", func() {
			logManagerMock.EXPECT().LogDebug(successfulPayloadGenerationMessage)
			payload, err := payloadGenerator.GeneratePayload(expirationTime)
			Expect(err).To(BeNil())
			Expect(payload).To(Not(BeNil()))
		})
	})
})
