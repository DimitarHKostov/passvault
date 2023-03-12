package session

var (
	sessionManager *SessionManager
)

type SessionManager struct {
	token string
}

func (s *SessionManager) Set(token string) {
	s.token = token
}

func (s *SessionManager) Get() string {
	return s.token
}

func (s *SessionManager) InvalidateSession() {
	s.token = ""
}

func Get() *SessionManager {
	if sessionManager == nil {
		sessionManager = &SessionManager{}
	}

	return sessionManager
}
