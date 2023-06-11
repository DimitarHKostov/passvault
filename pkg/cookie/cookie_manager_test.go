package cookie

// import (
// 	"passvault/pkg/jwt/mocks"
// 	"testing"

// 	logMock "passvault/pkg/log/mocks"

// 	"github.com/golang/mock/gomock"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )

// func TestCart(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "Cookie manager suite")
// }

// var _ = Describe("Shopping cart", func() {
// 	var (
// 		t              GinkgoTInterface
// 		jwtManagerMock *mocks.MockJWTManagerInterface
// 		logManagerMock *logMock.MockLogManagerInterface
// 		cookieManager  CookieManagerInterface
// 	)

// 	var _ = BeforeEach(func() {
// 		t = GinkgoT()
// 		jwtManagerMock = mocks.NewMockJWTManagerInterface(gomock.NewController(t))
// 		logManagerMock = logMock.NewMockLogManagerInterface(gomock.NewController(t))

// 	})

// 	Context("initially", func() {
// 		It("dummy test", func() {
// 			Expect(0).Should(BeZero())
// 		})
// 	})
// })
