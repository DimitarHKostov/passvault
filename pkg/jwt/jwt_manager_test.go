package jwt

import (
	"errors"
	"testing"
	"time"

	payloadGeneratorMock "passvault/pkg/generator/mocks"
	logMock "passvault/pkg/log/mocks"
	"passvault/pkg/types"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
		payloadMock    *payloadGeneratorMock.MockPayloadGeneratorInterface
		jwtManager     JWTManagerInterface
	)

	const (
		secretKey               = "b"
		expirationTime          = 10 * time.Minute
		expirationTimeInThePast = -10 * time.Minute
	)

	var _ = BeforeEach(func() {
		ctrl = gomock.NewController(ginkgo.GinkgoT())
		logManagerMock = logMock.NewMockLogManagerInterface(ctrl)
		payloadMock = payloadGeneratorMock.NewMockPayloadGeneratorInterface(ctrl)

		jwtManager = NewJwtManager(payloadMock, secretKey, logManagerMock)
	})

	var _ = AfterEach(func() {
		ctrl.Finish()
	})

	Context("GenerateToken", func() {
		It("should succeed to generate a new token", func() {
			payload := producePayload(expirationTime)
			payloadMock.EXPECT().GeneratePayload(expirationTime).Return(payload, nil)

			token, err := jwtManager.GenerateToken(expirationTime)
			Expect(err).To(BeNil())
			Expect(token).To(Not(Equal("")))
		})

		It("should fail to generate a new token when payload generator fails", func() {
			expectedError := errors.New("payload generation fails")
			payloadMock.EXPECT().GeneratePayload(expirationTime).Return(nil, expectedError)

			token, err := jwtManager.GenerateToken(expirationTime)
			Expect(err).To(Equal(expectedError))
			Expect(token).To(Equal(""))
		})
	})

	Context("VerifyToken", func() {
		It("should succeed to verify a token", func() {
			expectedPayload := producePayload(expirationTime)
			payloadMock.EXPECT().GeneratePayload(expirationTime).Return(expectedPayload, nil)
			token, err := jwtManager.GenerateToken(expirationTime)
			Expect(err).To(BeNil())
			Expect(token).To(Not(Equal("")))

			payload, err := jwtManager.VerifyToken(token)
			Expect(err).To(BeNil())
			Expect(payload.Uuid).To(Equal(expectedPayload.Uuid))
		})

		XIt("should return error when token is expired", func() {
			expectedPayload := producePayload(expirationTimeInThePast)
			payloadMock.EXPECT().GeneratePayload(expirationTimeInThePast).Return(expectedPayload, nil)
			token, err := jwtManager.GenerateToken(expirationTimeInThePast)
			Expect(err).To(BeNil())
			Expect(token).To(Not(Equal("")))

			payload, err := jwtManager.VerifyToken(token)
			Expect(err).To(BeNil())
			Expect(payload.Uuid).To(Equal(expectedPayload.Uuid))
		})
	})
})

func producePayload(duration time.Duration) *types.Payload {
	uuid, _ := uuid.NewRandom()

	payload := &types.Payload{
		Uuid:      uuid,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload
}
