package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/seanhuebl/unity-wealth/handlers"
	"github.com/seanhuebl/unity-wealth/internal/database"
	"github.com/stretchr/testify/assert"
)

type mockQueries struct {
	database.Queries
	CreateUserFunc func(ctx context.Context, params database.CreateUserParams) error
}

func (m *mockQueries) CreateUser(ctx context.Context, params database.CreateUserParams) error {
	return m.CreateUserFunc(ctx, params)
}

func TestAddUser(t *testing.T) {
	// Mock dependencies
	mockCfg := &handlers.ApiConfig{
		Queries: &mockQueries{},
	}

	router := gin.Default()
	router.POST("/addUser", func(ctx *gin.Context) {
		handlers.AddUser(ctx, mockCfg)
	})

	tests := []struct {
		name           string
		requestBody    string
		mockBehavior   func(*mockQueries)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Valid user creation",
			requestBody: `{"email":"user@example.com","password":"StrongPass123!"}`,
			mockBehavior: func(m *mockQueries) {
				m.CreateUserFunc = func(ctx context.Context, params database.CreateUserParams) error {
					assert.Equal(t, "user@example.com", params.Email)
					assert.NotEmpty(t, params.HashedPassword)
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Sign up successful!","email":"user@example.com"}`,
		},
		{
			name:           "Invalid JSON input",
			requestBody:    `{"email":"user@example.com","password":}`,
			mockBehavior:   func(m *mockQueries) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid character '}' looking for beginning of value"}`,
		},
		{
			name:        "Database error",
			requestBody: `{"email":"user@example.com","password":"StrongPass123!"}`,
			mockBehavior: func(m *mockQueries) {
				m.CreateUserFunc = func(ctx context.Context, params database.CreateUserParams) error {
					return errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockBehavior(mockCfg.Queries.(*mockQueries))

			req := httptest.NewRequest(http.MethodPost, "/addUser", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
