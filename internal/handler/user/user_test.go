package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"

	"github.com/Avito-courses/l11-examples/internal/handler/user/mocks"
	model "github.com/Avito-courses/l11-examples/internal/model/user"
)

func TestController_Get(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockResponse   *model.User
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "success",
			userID: "1",
			mockResponse: &model.User{
				ID:    1,
				Name:  "John Doe",
				Phone: "123123",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":    float64(1),
				"name":  "John Doe",
				"phone": "123123",
			},
		},
		{
			name:           "Пользователь не найден",
			userID:         "999",
			mockError:      model.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "User not found",
			},
		},
		{
			name:           "Внутренняя ошибка сервера",
			userID:         "2",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Internal server error",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockuserRepo(ctrl)
			mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tc.mockResponse, tc.mockError)

			controller := &Controller{
				repo: mockRepo,
			}

			req := httptest.NewRequest("GET", "/users/"+tc.userID, nil)

			rr := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/users/{id}", controller.Get)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			controller.Get(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want application/json",
					contentType)
			}

			var response map[string]interface{}
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			for key, expectedValue := range tc.expectedBody {
				if value, ok := response[key]; !ok || value != expectedValue {
					t.Errorf("handler returned unexpected %s: got %v want %v",
						key, value, expectedValue)
				}
			}
		})
	}
}
