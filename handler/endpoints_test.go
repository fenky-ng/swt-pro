package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	errorHelper "github.com/fenky-ng/swt-pro/helper/error"
	"github.com/fenky-ng/swt-pro/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func Test_Server_Register(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mock           func(fields *fields)
		wantStatusCode int
		wantErr        error
	}{
		{
			name: "no request body",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(``)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "invalid request",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "",
						"full_name": "",
						"password": ""
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "error checkNewPhoneNumber",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"full_name": "Sawit Pro 1",
						"password": "Sawit@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{}, errors.New("expected GetUserByPhoneNumber error")).
					Times(1)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        nil,
		},
		{
			name: "phone number already registered",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"full_name": "Sawit Pro 1",
						"password": "Sawit@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{
						ID: 1,
					}, nil).
					Times(1)
			},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "error InsertUser",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"full_name": "Sawit Pro 1",
						"password": "Sawit@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{
						ID: 0,
					}, nil).
					Times(1)

				fields.Repository.EXPECT().InsertUser(context.Background(), gomock.AssignableToTypeOf(repository.User{})).
					Return(int64(0), errors.New("expected InsertUser error")).
					Times(1)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        nil,
		},
		{
			name: "passed",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"full_name": "Sawit Pro 1",
						"password": "Sawit@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{
						ID: 0,
					}, nil).
					Times(1)

				fields.Repository.EXPECT().InsertUser(context.Background(), gomock.AssignableToTypeOf(repository.User{})).
					Return(int64(1), nil).
					Times(1)
			},
			wantStatusCode: http.StatusOK,
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			gotErr := s.Register(tt.args.ctx)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("Server.Register() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if gotErr == nil {
				if tt.args.ctx.Response().Status != tt.wantStatusCode {
					t.Errorf("Server.Register() gotStatusCode = %d, wantStatusCode = %d", tt.args.ctx.Response().Status, tt.wantStatusCode)
				}
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}

func Test_Server_Login(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mock           func(fields *fields)
		wantStatusCode int
		wantErr        error
	}{
		{
			name: "no request body",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(``)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "error GetUserByPhoneNumber",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"password": "Sawit@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{}, errors.New("expected GetUserByPhoneNumber error")).
					Times(1)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        nil,
		},
		{
			name: "user not found",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"password": "Sawit@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{
						ID: 0,
					}, nil).
					Times(1)
			},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "wrong password",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"password": "Random@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{
						ID:       1,
						Password: "$2a$04$DcEZFEpGx1t/cpN1jHBjrO2wRLM317fSp.aU4uQtw3GUhbDMvXODe",
					}, nil).
					Times(1)
			},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "error IncreaseLoginCount",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"password": "Sawit@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{
						ID:       1,
						Password: "$2a$04$DcEZFEpGx1t/cpN1jHBjrO2wRLM317fSp.aU4uQtw3GUhbDMvXODe",
					}, nil).
					Times(1)

				fields.Repository.EXPECT().IncreaseLoginCount(context.Background(), int64(1)).
					Return(errors.New("expected IncreaseLoginCount error")).
					Times(1)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        nil,
		},
		{
			name: "passed",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"password": "Sawit@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{
						ID:       1,
						Password: "$2a$04$DcEZFEpGx1t/cpN1jHBjrO2wRLM317fSp.aU4uQtw3GUhbDMvXODe",
					}, nil).
					Times(1)

				fields.Repository.EXPECT().IncreaseLoginCount(context.Background(), int64(1)).
					Return(nil).
					Times(1)
			},
			wantStatusCode: http.StatusOK,
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			gotErr := s.Login(tt.args.ctx)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("Server.Login() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if gotErr == nil {
				if tt.args.ctx.Response().Status != tt.wantStatusCode {
					t.Errorf("Server.Register() gotStatusCode = %d, wantStatusCode = %d", tt.args.ctx.Response().Status, tt.wantStatusCode)
				}
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}

func Test_Server_GetProfile(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mock           func(fields *fields)
		wantStatusCode int
		wantErr        error
	}{
		{
			name: "invalid authorization",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusForbidden,
			wantErr:        nil,
		},
		{
			name: "error GetUserByID",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByID(context.Background(), int64(1)).
					Return(repository.User{}, errors.New("expected GetUserByID error")).
					Times(1)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        nil,
		},
		{
			name: "passed",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByID(context.Background(), int64(1)).
					Return(repository.User{}, nil).
					Times(1)
			},
			wantStatusCode: http.StatusOK,
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			gotErr := s.GetProfile(tt.args.ctx)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("Server.GetProfile() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if gotErr == nil {
				if tt.args.ctx.Response().Status != tt.wantStatusCode {
					t.Errorf("Server.Register() gotStatusCode = %d, wantStatusCode = %d", tt.args.ctx.Response().Status, tt.wantStatusCode)
				}
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}

func Test_Server_UpdateProfile(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mock           func(fields *fields)
		wantStatusCode int
		wantErr        error
	}{
		{
			name: "invalid authorization",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", nil)
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusForbidden,
			wantErr:        nil,
		},
		{
			name: "no request body",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(``)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "no params at all",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{}`)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "all params have empty value",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "",
						"full_name": ""
					}`)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "all params have empty value",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "",
						"full_name": ""
					}`)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "invalid phone number",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "123",
						"full_name": "Sawit Pro 1"
					}`)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "error GetUserByPhoneNumber",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"full_name": "Sawit Pro 1"
					}`)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{}, errors.New("expected GetUserByPhoneNumber error")).
					Times(1)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        nil,
		},
		{
			name: "phone number is already registered",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"full_name": "Sawit Pro 1"
					}`)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{
						ID: 2,
					}, nil).
					Times(1)
			},
			wantStatusCode: http.StatusConflict,
			wantErr:        nil,
		},
		{
			name: "invalid full name",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"full_name": "SP"
					}`)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{}, nil).
					Times(1)
			},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        nil,
		},
		{
			name: "error UpdateUser",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"full_name": "Sawit Pro 1"
					}`)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{}, nil).
					Times(1)

				fields.Repository.EXPECT().UpdateUser(context.Background(),
					repository.User{
						ID:          1,
						PhoneNumber: "+628223344551",
						FullName:    "Sawit Pro 1",
					}).
					Return(errors.New("expected UpdateUser error")).
					Times(1)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        nil,
		},
		{
			name: "passed",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+628223344551",
						"full_name": "Sawit Pro 1"
					}`)))
					jwt, _ := generateJwtToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344551").
					Return(repository.User{}, nil).
					Times(1)

				fields.Repository.EXPECT().UpdateUser(context.Background(),
					repository.User{
						ID:          1,
						PhoneNumber: "+628223344551",
						FullName:    "Sawit Pro 1",
					}).
					Return(nil).
					Times(1)
			},
			wantStatusCode: http.StatusOK,
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			gotErr := s.UpdateProfile(tt.args.ctx)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("Server.UpdateProfile() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if gotErr == nil {
				if tt.args.ctx.Response().Status != tt.wantStatusCode {
					t.Errorf("Server.Register() gotStatusCode = %d, wantStatusCode = %d", tt.args.ctx.Response().Status, tt.wantStatusCode)
				}
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}
