package cookie

import (
	"errors"
	"net/http"
	"passvault/pkg/jwt/mocks"
	"passvault/pkg/types"
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
		ctrl           *gomock.Controller
		jwtManagerMock *mocks.MockJWTManagerInterface
		logManagerMock *logMock.MockLogManagerInterface
		cookieManager  CookieManagerInterface
	)

	const (
		expirationTime = 10 * time.Minute
		tokenValue     = "token-value"
		errorMessage   = "error occurred while creating token"
	)

	var _ = BeforeEach(func() {
		ctrl = gomock.NewController(ginkgo.GinkgoT())
		jwtManagerMock = mocks.NewMockJWTManagerInterface(ctrl)
		logManagerMock = logMock.NewMockLogManagerInterface(ctrl)

		cookieManager = NewCookieManager(jwtManagerMock, logManagerMock)
	})

	var _ = AfterEach(func() {
		ctrl.Finish()
	})

	Context("ProduceCookie", func() {
		It("should succeed", func() {
			expectedCookie := produceMockCookie(types.CookieName, tokenValue, expirationTime, types.CookieHttpOnly)
			jwtManagerMock.EXPECT().GenerateToken(expirationTime).Return(tokenValue, nil)
			actualCookie, err := cookieManager.ProduceCookie()

			Expect(err).To(BeNil())
			Expect(actualCookie.Name).To(Equal(expectedCookie.Name))
			Expect(actualCookie.Value).To(Equal(expectedCookie.Value))
			Expect(actualCookie.HttpOnly).To(Equal(expectedCookie.HttpOnly))
		})

		It("should fail", func() {
			expectedError := errors.New(errorMessage)
			jwtManagerMock.EXPECT().GenerateToken(expirationTime).Return("", expectedError)
			actualCookie, err := cookieManager.ProduceCookie()

			Expect(err).To(Equal(expectedError))
			Expect(actualCookie).To(BeNil())
		})
	})
})

func produceMockCookie(cookieName, token string, expirationTime time.Duration, httpOnly bool) http.Cookie {
	return http.Cookie{Name: cookieName, Value: token, Expires: time.Now().Add(expirationTime), HttpOnly: httpOnly}
}
