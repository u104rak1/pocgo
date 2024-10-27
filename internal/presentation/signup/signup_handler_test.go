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
		requestBody      interface{}
		prepare          func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase)
		expectedCode     int
		expectedReason   string
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
			expectedCode: http.StatusCreated,
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
			name:         "Error occurs during signup when request body is invalid json.",
			requestBody:  "invalid json",
			prepare:      func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedResponse: response.ErrorResponse{
				Reason:  response.BadRequestReason,
				Message: "code=400, message=Unmarshal type error: expected=signup.SignupRequestBody, got=string, field=, offset=14, internal=json: cannot unmarshal string into Go value of type signup.SignupRequestBody",
			},
		},
		{
			name:           "Error occurs during signup when validation failed.",
			requestBody:    signup.SignupRequestBody{},
			prepare:        func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {},
			expectedCode:   http.StatusBadRequest,
			expectedReason: response.ValidationFailedReason,
		},
		{
			name:        "Error occurs during signup when user email already exists.",
			requestBody: validRequestBody,
			prepare: func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, userDomain.ErrUserEmailAlreadyExists)
			},
			expectedCode: http.StatusConflict,
			expectedResponse: response.ErrorResponse{
				Reason:  response.ConflictReason,
				Message: userDomain.ErrUserEmailAlreadyExists.Error(),
			},
		},
		{
			name:        "Error occurs during signup when authentication already exists.",
			requestBody: validRequestBody,
			prepare: func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, authDomain.ErrAuthenticationAlreadyExists)
			},
			expectedCode: http.StatusConflict,
			expectedResponse: response.ErrorResponse{
				Reason:  response.ConflictReason,
				Message: authDomain.ErrAuthenticationAlreadyExists.Error(),
			},
		},
		{
			name:        "Error occurs during signup when unknown error occurs.",
			requestBody: validRequestBody,
			prepare: func(ctx context.Context, mockSignupUC *authMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, errors.New("unknown error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResponse: response.ErrorResponse{
				Reason:  response.InternalServerErrorReason,
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
			assert.Equal(t, tt.expectedCode, rec.Code)
			if rec.Code == http.StatusCreated {
				var resp signup.SignupResponseBody
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, resp)
			} else if rec.Code == http.StatusBadRequest && tt.expectedReason == response.ValidationFailedReason {
				var resp response.ValidationErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, response.ValidationFailedReason, resp.Reason)
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
