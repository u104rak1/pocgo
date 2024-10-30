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
	appMock "github.com/ucho456job/pocgo/internal/application/mock"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
	"github.com/ucho456job/pocgo/internal/presentation/signup"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestSignupHandler(t *testing.T) {
	var (
		userID           = ulid.GenerateStaticULID("user")
		userName         = "sato taro"
		userEmail        = "sato@example.com"
		userPassword     = "password"
		accountID        = ulid.GenerateStaticULID("account")
		accountName      = "For work"
		accountBalance   = 0.0
		accountPassword  = "1234"
		accountCurrency  = money.JPY
		accountUpdatedAt = "2023-10-20T00:00:00Z"
		accessToken      = "token"
		unknownErr       = errors.New("unknown error")
	)

	var requestBody = signup.SignupRequestBody{
		User: signup.SignupRequestBodyUser{
			Name:     userName,
			Email:    userEmail,
			Password: userPassword,
			Account: signup.SignupRequestBodyAccount{
				Name:     accountName,
				Password: accountPassword,
				Currency: accountCurrency,
			},
		},
	}

	tests := []struct {
		caseName             string
		requestBody          interface{}
		prepare              func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase)
		expectedCode         int
		expectedReason       string
		expectedResponseBody interface{}
	}{
		{
			caseName:    "Successful signup.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(&authApp.SignupDTO{
					User: userApp.CreateUserDTO{
						ID:    userID,
						Name:  userName,
						Email: userEmail,
					},
					Account: accountApp.CreateAccountDTO{
						ID:        accountID,
						Name:      accountName,
						Balance:   accountBalance,
						Currency:  accountCurrency,
						UpdatedAt: accountUpdatedAt,
					},
					AccessToken: accessToken,
				}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedResponseBody: signup.SignupResponseBody{
				User: signup.SignupResponseBodyUser{
					ID:    userID,
					Name:  userName,
					Email: userEmail,
					Account: signup.SignupResponseBodyAccount{
						ID:        accountID,
						Name:      accountName,
						Balance:   accountBalance,
						Currency:  accountCurrency,
						UpdatedAt: accountUpdatedAt,
					},
				},
				AccessToken: accessToken,
			},
		},
		{
			caseName:     "Error occurs during signup when request body is invalid json.",
			requestBody:  "invalid json",
			prepare:      func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {},
			expectedCode: http.StatusBadRequest,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.BadRequestReason,
				Message: "code=400, message=Unmarshal type error: expected=signup.SignupRequestBody, got=string, field=, offset=14, internal=json: cannot unmarshal string into Go value of type signup.SignupRequestBody",
			},
		},
		{
			caseName:       "Error occurs during signup when validation failed.",
			requestBody:    signup.SignupRequestBody{},
			prepare:        func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {},
			expectedCode:   http.StatusBadRequest,
			expectedReason: response.ValidationFailedReason,
		},
		{
			caseName:    "Error occurs during signup when user email already exists.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, userDomain.ErrEmailAlreadyExists)
			},
			expectedCode: http.StatusConflict,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.ConflictReason,
				Message: userDomain.ErrEmailAlreadyExists.Error(),
			},
		},
		{
			caseName:    "Error occurs during signup when authentication already exists.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, authDomain.ErrAlreadyExists)
			},
			expectedCode: http.StatusConflict,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.ConflictReason,
				Message: authDomain.ErrAlreadyExists.Error(),
			},
		},
		{
			caseName:    "Unknown error occurs during signup.",
			requestBody: requestBody,
			prepare: func(ctx context.Context, mockSignupUC *appMock.MockISignupUsecase) {
				mockSignupUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, unknownErr)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResponseBody: response.ErrorResponse{
				Reason:  response.InternalServerErrorReason,
				Message: unknownErr.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			mockSignupUC := appMock.NewMockISignupUsecase(ctrl)
			tt.prepare(ctx.Request().Context(), mockSignupUC)

			h := signup.NewSignupHandler(mockSignupUC)
			err := h.Run(ctx)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if rec.Code == http.StatusCreated {
				var resp signup.SignupResponseBody
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseBody, resp)
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
				assert.Equal(t, tt.expectedResponseBody, resp)
			}
		})
	}
}
