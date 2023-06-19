package app

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cookieMock "passvault/pkg/cookie/mocks"
	cryptMock "passvault/pkg/crypt/mocks"
	databaseMock "passvault/pkg/database/mocks"
	logMock "passvault/pkg/log/mocks"
	middleMock "passvault/pkg/middleware/mocks"
	"passvault/pkg/types"

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
		ctrl                *gomock.Controller
		logManagerMock      *logMock.MockLogManagerInterface
		cookieManagerMock   *cookieMock.MockCookieManagerInterface
		cryptManagerMock    *cryptMock.MockCryptManagerInterface
		databaseManagerMock *databaseMock.MockDatabaseManagerInterface
		middlewareMock      *middleMock.MockMiddlewareInterface
		app                 *App

		responseRecorder httptest.ResponseRecorder
	)

	const (
		mockExpirationTime  = 5 * time.Minute
		mockToken           = "token-value"
		mockDomain          = "domain.com"
		mockShortUsername   = "short"
		mockShortPassword   = "short"
		mockCorrectUsername = "correct@test.com"
		mockCorrectPassword = "correctPassword98765"
	)

	var _ = BeforeEach(func() {
		ctrl = gomock.NewController(ginkgo.GinkgoT())
		logManagerMock = logMock.NewMockLogManagerInterface(ctrl)
		cookieManagerMock = cookieMock.NewMockCookieManagerInterface(ctrl)
		cryptManagerMock = cryptMock.NewMockCryptManagerInterface(ctrl)
		databaseManagerMock = databaseMock.NewMockDatabaseManagerInterface(ctrl)
		middlewareMock = middleMock.NewMockMiddlewareInterface(ctrl)

		app = NewApp()
		app.AppOpts.LogManager = logManagerMock
		app.AppOpts.CookieManager = cookieManagerMock
		app.AppOpts.CryptManager = cryptManagerMock
		app.AppOpts.DatabaseManager = databaseManagerMock
		app.AppOpts.Middleware = middlewareMock

		responseRecorder = *httptest.NewRecorder()
	})

	var _ = AfterEach(func() {
		ctrl.Finish()
	})

	Context("Login", func() {
		It("should fail with 401 when the body in the request is empty", func() {
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(""))
			Expect(err).To(BeNil())
			logManagerMock.EXPECT().LogDebug(types.EmptyBodyMessage)

			app.login(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusUnauthorized))
		})

		It("should fail with 400 when the credentials in the body are not in json format", func() {
			request, err := produceRequest(http.MethodGet, "", strings.NewReader("<html>test</html>"))
			Expect(err).To(BeNil())
			logManagerMock.EXPECT().LogError("invalid character '<' looking for beginning of value")

			app.login(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should fail with 401 when the credentials(password) in the body are wrong", func() {
			wrongCredentials := produceCredentials("wrong-password")
			wrongCredentialsJson, err := json.Marshal(&wrongCredentials)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(wrongCredentialsJson)))
			Expect(err).To(BeNil())
			logManagerMock.EXPECT().LogDebug("crypto/bcrypt: hashedPassword is not the hash of the given password")

			app.login(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusUnauthorized))
		})

		It("should fail with 500 when the cookie manager returns error", func() {
			correctCredentials := produceCredentials("VeryHard987*") //todo refactor
			correctCredentialsJson, err := json.Marshal(&correctCredentials)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(correctCredentialsJson)))
			Expect(err).To(BeNil())
			cookieManagerMock.EXPECT().ProduceCookie().Return(nil, errors.New("error"))
			logManagerMock.EXPECT().LogError("error")

			app.login(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should succeed with 200 and return cookie", func() {
			correctCredentials := produceCredentials("VeryHard987*") //todo refactor
			correctCredentialsJson, err := json.Marshal(&correctCredentials)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(correctCredentialsJson)))
			Expect(err).To(BeNil())
			cookie := produceCookie(mockExpirationTime, mockToken)
			cookieManagerMock.EXPECT().ProduceCookie().Return(cookie, nil)
			logManagerMock.EXPECT().LogDebug(successfulLoginMessage)

			app.login(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusOK))
			Expect(responseRecorder.Result().Header.Get("Set-Cookie")).To(Equal(cookie.String()))
		})
	})

	Context("Save", func() {
		It("should fail with 400 when the body in the request is empty", func() {
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(""))
			Expect(err).To(BeNil())
			logManagerMock.EXPECT().LogDebug(types.EmptyBodyMessage)

			app.save(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should fail with 400 when the entry in the body is not in json format", func() {
			request, err := produceRequest(http.MethodGet, "", strings.NewReader("<html>test</html>"))
			Expect(err).To(BeNil())
			logManagerMock.EXPECT().LogError("invalid character '<' looking for beginning of value")

			app.save(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should fail with 400 when the username in the entry is less than 6 symbols", func() {
			entry := produceEntry(mockDomain, mockShortUsername, mockShortPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			logManagerMock.EXPECT().LogDebug("input too short, must be at least 6 characters")

			app.save(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should fail with 400 when the password in the entry is less than 6 symbols", func() {
			entry := produceEntry(mockDomain, mockShortUsername, mockShortPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			logManagerMock.EXPECT().LogDebug("input too short, must be at least 6 characters")

			app.save(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should fail with 500 when the database manager fails to return if domain already exists", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(false, errors.New("error"))
			logManagerMock.EXPECT().LogError("error")

			app.save(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should fail with 409 when the domain already exists", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(true, nil)
			logManagerMock.EXPECT().LogDebug(domainAlreadyExistsMessage)

			app.save(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusConflict))
		})

		It("should fail with 500 when the crypt manager fails", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(false, nil)
			cryptManagerMock.EXPECT().Encrypt(mockCorrectPassword).Return(nil, errors.New("error"))
			logManagerMock.EXPECT().LogError("error")

			app.save(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should fail with 500 when the database manager fails to save to the database", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(false, nil)
			correctPasswordPtr := new(string)
			*correctPasswordPtr = mockCorrectPassword
			cryptManagerMock.EXPECT().Encrypt(mockCorrectPassword).Return(correctPasswordPtr, nil)
			databaseManagerMock.EXPECT().Save(entry).Return(errors.New("error"))
			logManagerMock.EXPECT().LogError("error")

			app.save(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should return 201 and succeed to save entry to database", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(false, nil)
			correctPasswordPtr := new(string)
			*correctPasswordPtr = mockCorrectPassword
			cryptManagerMock.EXPECT().Encrypt(mockCorrectPassword).Return(correctPasswordPtr, nil)
			databaseManagerMock.EXPECT().Save(entry).Return(nil)
			logManagerMock.EXPECT().LogDebug(successfulSaveMessage)

			app.save(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusCreated))
		})
	})

	Context("Retrieve", func() {
		It("should fail with 400 when the body in the request is empty", func() {
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(""))
			Expect(err).To(BeNil())
			logManagerMock.EXPECT().LogDebug(types.EmptyBodyMessage)

			app.retrieve(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should fail with 400 when the entry in the body is not in json format", func() {
			request, err := produceRequest(http.MethodGet, "", strings.NewReader("<html>test</html>"))
			Expect(err).To(BeNil())
			logManagerMock.EXPECT().LogError("invalid character '<' looking for beginning of value")

			app.retrieve(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
		})

		XIt("should fail with 400 when the domain validation fails", func() {
			entry := produceEntry(mockDomain, "", "")
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())

			logManagerMock.EXPECT().LogDebug("msg")

			app.retrieve(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should fail with 500 when the database manager fails to return if domain already exists", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(false, errors.New("error"))
			logManagerMock.EXPECT().LogError("error")

			app.retrieve(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should fail with 404 when the domain does not exists", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(false, nil)
			logManagerMock.EXPECT().LogDebug(domainDoesNotExistMessage)

			app.retrieve(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should fail with 500 when the database manager fails to return the corresponding entry", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(true, nil)
			databaseManagerMock.EXPECT().Get(mockDomain).Return(nil, errors.New("error"))
			logManagerMock.EXPECT().LogError("error")

			app.retrieve(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should fail with 500 when the crypt manager fails to decrypt the password", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(true, nil)
			databaseManagerMock.EXPECT().Get(mockDomain).Return(&entry, nil) //tuka
			cryptManagerMock.EXPECT().Decrypt(mockCorrectPassword).Return(nil, errors.New("error"))
			logManagerMock.EXPECT().LogError("error")

			app.retrieve(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("succeed with 200 and return the entry from the db", func() {
			entry := produceEntry(mockDomain, mockCorrectUsername, mockCorrectPassword)
			entryJson, err := json.Marshal(entry)
			Expect(err).To(BeNil())
			request, err := produceRequest(http.MethodGet, "", strings.NewReader(string(entryJson)))
			Expect(err).To(BeNil())
			databaseManagerMock.EXPECT().Contains(mockDomain).Return(true, nil)
			databaseManagerMock.EXPECT().Get(mockDomain).Return(&entry, nil)
			correctPasswordPtr := new(string)
			*correctPasswordPtr = mockCorrectPassword
			cryptManagerMock.EXPECT().Decrypt(mockCorrectPassword).Return(correctPasswordPtr, nil)
			logManagerMock.EXPECT().LogDebug(successfulRetrieveMessage)

			app.retrieve(&responseRecorder, request)
			Expect(responseRecorder.Result().StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("Update", func() {

	})
})

func produceRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func produceCredentials(password string) types.Credentials {
	return types.Credentials{Password: password}
}

func produceCookie(expirationTime time.Duration, token string) *http.Cookie {
	return &http.Cookie{Name: types.CookieName, Value: token, Expires: time.Now().Add(expirationTime), HttpOnly: types.CookieHttpOnly}
}

func produceEntry(domain, username, password string) types.Entry {
	return types.Entry{Domain: domain, Username: username, Password: password}
}
