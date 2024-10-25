package signup_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	authApp "github.com/ucho456job/pocgo/internal/application/authentication"
	authMock "github.com/ucho456job/pocgo/internal/application/mock"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
	"github.com/ucho456job/pocgo/internal/presentation/signup"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestSignupHandler_Run(t *testing.T) {
	var (
		validUserID           = ulid.GenerateStaticULID("user")
		validUserName         = "sato taro"
		validUserEmail        = "sato@example.com"
		validUserPassword     = "password"
		validAccountID        = ulid.GenerateStaticULID("account")
		validAccountName      = "For work"
		validAccountBalance   = 0.0
		validAccountPassword  = "1234"
		validAccountCurrency  = money.JPY
		validAccountUpdatedAt = "2023-10-20T00:00:00Z"
		validAccessToken      = "token"
	)

	var validRequestBody = signup.SignupRequestBody{
		User: signup.SignupRequestBodyUser{
			Name:     validUserName,
			Email:    validUserEmail,
			Password: validUserPassword,
			Account: signup.SignupRequestBodyAccount{
				Name:     validAccountName,
				Password: validAccountPassword,
				Currency: validAccountCurrency,
			},
		},
	}

	tests := []struct {
		name             string
		requestBody      signup.SignupRequestBody
		prepare          func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:        "Successful signup.",
			requestBody: validRequestBody,
			prepare: func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(&authApp.SignupDTO{
					User: userApp.CreateUserDTO{
						ID:    validUserID,
						Name:  validUserName,
						Email: validUserEmail,
					},
					Account: accountApp.CreateAccountDTO{
						ID:        validAccountID,
						Name:      validAccountName,
						Balance:   validAccountBalance,
						Currency:  validAccountCurrency,
						UpdatedAt: validAccountUpdatedAt,
					},
					AccessToken: validAccessToken,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedResponse: signup.SignupResponseBody{
				User: signup.SignupResponseBodyUser{
					ID:    validUserID,
					Name:  validUserName,
					Email: validUserEmail,
					Account: signup.SignupResponseBodyAccount{
						ID:        validAccountID,
						Name:      validAccountName,
						Balance:   validAccountBalance,
						Currency:  validAccountCurrency,
						UpdatedAt: validAccountUpdatedAt,
					},
				},
				AccessToken: validAccessToken,
			},
		},
		{
			name:           "Error occurs during signup when request body is invalid.",
			requestBody:    signup.SignupRequestBody{},
			prepare:        func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Error occurs during signup when user email already exists.",
			requestBody: validRequestBody,
			prepare: func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, userDomain.ErrUserEmailAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedResponse: response.ErrorResponse{
				Code:    response.ConflictCode,
				Message: userDomain.ErrUserEmailAlreadyExists.Error(),
			},
		},
		{
			name:        "Error occurs during signup when authentication already exists.",
			requestBody: validRequestBody,
			prepare: func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, authDomain.ErrAuthenticationAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedResponse: response.ErrorResponse{
				Code:    response.ConflictCode,
				Message: authDomain.ErrAuthenticationAlreadyExists.Error(),
			},
		},
		{
			name:        "Error occurs during signup when unknown error occurs.",
			requestBody: validRequestBody,
			prepare: func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, errors.New("unknown error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: response.ErrorResponse{
				Code:    response.InternalServerErrorCode,
				Message: "unknown error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			mockSignupUC := authMock.NewMockISignupUsecase(ctrl)
			tt.prepare(ctx.Request().Context(), mockSignupUC)
			h := signup.NewSignupHandler(mockSignupUC)

			err := h.Run(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if rec.Code == http.StatusCreated {
				var resp signup.SignupResponseBody
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, resp)
			} else if rec.Code == http.StatusBadRequest {
				var resp response.ValidationErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, response.ValidationFailedCode, resp.Code)
				assert.NotEmpty(t, resp.Errors)
			} else {
				var resp response.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, resp)
			}
		})
	}
}
