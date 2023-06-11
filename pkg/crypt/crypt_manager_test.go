package crypt

import (
	"crypto/aes"
	"net/http"
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
		logManagerMock *logMock.MockLogManagerInterface
		cryptManager   CryptManagerInterface
	)

	const (
		badCryptManagerSecretKey  = "bad-secret-key"
		goodCryptManagerSecretKey = "this is secret key enough 32 bit"
		testPassword              = "test-password"
		testHex                   = "1234f"
	)

	var _ = BeforeEach(func() {
		ctrl = gomock.NewController(ginkgo.GinkgoT())
		logManagerMock = logMock.NewMockLogManagerInterface(ctrl)
	})

	var _ = AfterEach(func() {
		ctrl.Finish()
	})

	Context("Encrypt", func() {
		It("should succeed to encrypt the password", func() {
			cryptManager = NewCryptManager(logManagerMock, []byte(goodCryptManagerSecretKey))

			encToString, err := cryptManager.Encrypt(testPassword)
			Expect(err).To(BeNil())
			Expect(encToString).To(Not(BeNil()))
		})

		It("should fail to encrypt the password when the given crypter secret key is not 32 bits", func() {
			cryptManager = NewCryptManager(logManagerMock, []byte(badCryptManagerSecretKey))

			expectedError := aes.KeySizeError(len(badCryptManagerSecretKey))
			encToString, err := cryptManager.Encrypt(testPassword)

			Expect(err).To(Equal(expectedError))
			Expect(encToString).To(BeNil())
		})
	})

	Context("Decrypt", func() {
		It("should succeed to decrypt the password", func() {
			cryptManager = NewCryptManager(logManagerMock, []byte(goodCryptManagerSecretKey))

			encToString, err := cryptManager.Encrypt(testPassword)
			Expect(err).To(BeNil())

			decrypted, err := cryptManager.Decrypt(*encToString)
			Expect(err).To(BeNil())

			Expect(decrypted).To(Not(BeNil()))
			Expect(*decrypted).To(Equal(testPassword))
		})

		It("should fail to decrypt the password when the given crypter secret key is not 32 bits", func() {
			cryptManager = NewCryptManager(logManagerMock, []byte(badCryptManagerSecretKey))

			expectedError := aes.KeySizeError(len(badCryptManagerSecretKey))
			decrypted, err := cryptManager.Decrypt(testHex)

			Expect(err).To(Equal(expectedError))
			Expect(decrypted).To(BeNil())
		})
	})
})

func produceMockCookie(cookieName, token string, expirationTime time.Duration, httpOnly bool) http.Cookie {
	return http.Cookie{Name: cookieName, Value: token, Expires: time.Now().Add(expirationTime), HttpOnly: httpOnly}
}
