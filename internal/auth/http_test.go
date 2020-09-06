package auth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (r *UserRepositoryMock) Create(user *User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *UserRepositoryMock) GetUserByName(name string) (*User, error) {
	args := r.Called(name)
	return args.Get(0).(*User), args.Error(1)
}

func TestCreateUserHandler(t *testing.T) {
	w := httptest.NewRecorder()
	repo := &UserRepositoryMock{}
	repo.On("Create", mock.Anything).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "whatever", strings.NewReader(`{"Name": "thomas", "password": "PSG2021"}`))
	CreateUserHandler(repo).ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "thomas", gjson.Get(string(body), "Name").Str)
	assert.True(t, gjson.Get(string(body), "TeamID").Exists())
	repo.AssertExpectations(t)
}

func TestAuthMiddlewareWithCredentials(t *testing.T) {
	nextHandlerCalled := false
	var currentUser *User
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
		currentUser = CurrentUserFromContext(r.Context())
	})

	w := httptest.NewRecorder()
	repo := &UserRepositoryMock{}
	hash, _ := HashPassword("password")
	user := &User{
		Name:         "thomas",
		PasswordHash: hash,
		TeamID:       uuid.NewV4(),
	}
	repo.On("GetUserByName", "thomas").Return(user, nil)

	req, _ := http.NewRequest(http.MethodGet, "whatever", nil)
	req.SetBasicAuth("thomas", "password")

	AuthMiddleware(repo)(handler).ServeHTTP(w, req)

	assert.True(t, nextHandlerCalled)
	assert.Equal(t, user, currentUser)
	repo.AssertExpectations(t)
}

func TestAuthMiddlewareWithoutCredentials(t *testing.T) {
	nextHandlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
	})

	w := httptest.NewRecorder()
	repo := &UserRepositoryMock{}

	req, _ := http.NewRequest(http.MethodGet, "whatever", nil)

	AuthMiddleware(repo)(handler).ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.False(t, nextHandlerCalled)
	repo.AssertExpectations(t)
}
