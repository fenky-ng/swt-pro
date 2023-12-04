package handler

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/fenky-ng/swt-pro/generated"
	errorHelper "github.com/fenky-ng/swt-pro/helper/error"
	"github.com/fenky-ng/swt-pro/repository"
	"github.com/golang/mock/gomock"
)

func Test_validateRegistration(t *testing.T) {
	type args struct {
		request generated.RegistrationRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
	}{
		{
			name: "invalid phone number",
			args: args{
				request: generated.RegistrationRequest{
					PhoneNumber: "+62812345",
					FullName:    "Sawit Pro",
					Password:    "Sawit@Pr0",
				},
			},
			wantRes: []string{
				"Phone numbers must be at minimum 10 characters and maximum 13 characters",
			},
		},
		{
			name: "invalid full name",
			args: args{
				request: generated.RegistrationRequest{
					PhoneNumber: "+628123456",
					FullName:    "SP",
					Password:    "Sawit@Pr0",
				},
			},
			wantRes: []string{
				"Full name must be at minimum 3 characters and maximum 60 characters",
			},
		},
		{
			name: "invalid password",
			args: args{
				request: generated.RegistrationRequest{
					PhoneNumber: "+628123456",
					FullName:    "Sawit Pro",
					Password:    "sawitpro",
				},
			},
			wantRes: []string{
				"Passwords must be minimum 6 characters and maximum 64 characters, containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := validateRegistration(tt.args.request)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("validateRegistration() gotRes = %+v, wantRes = %+v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_validatePhoneNumber(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
	}{
		{
			name: "too short",
			args: args{
				input: "+62822334",
			},
			wantRes: []string{
				"Phone numbers must be at minimum 10 characters and maximum 13 characters",
			},
		},
		{
			name: "too long",
			args: args{
				input: "+6282233445566",
			},
			wantRes: []string{
				"Phone numbers must be at minimum 10 characters and maximum 13 characters",
			},
		},
		{
			name: "invalid prefix",
			args: args{
				input: "1234567890",
			},
			wantRes: []string{
				"Phone numbers must start with the Indonesia country code “+62”",
			},
		},
		{
			name: "passed",
			args: args{
				input: "+628223344556",
			},
			wantRes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := validatePhoneNumber(tt.args.input)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("validatePhoneNumber() gotRes = %+v, wantRes = %+v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_validateFullName(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
	}{
		{
			name: "too short",
			args: args{
				input: "12",
			},
			wantRes: []string{
				"Full name must be at minimum 3 characters and maximum 60 characters",
			},
		},
		{
			name: "too long",
			args: args{
				input: "1234567890" +
					"1234567890" +
					"1234567890" +
					"1234567890" +
					"1234567890" +
					"1234567890" +
					"1",
			},
			wantRes: []string{
				"Full name must be at minimum 3 characters and maximum 60 characters",
			},
		},
		{
			name: "passed",
			args: args{
				input: "1234567890",
			},
			wantRes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := validateFullName(tt.args.input)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("validateFullName() gotRes = %+v, wantRes = %+v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_validatePassword(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		wantRes bool
	}{
		{
			name: "too short",
			args: args{
				input: "12",
			},
			wantRes: false,
		},
		{
			name: "too long",
			args: args{
				input: "12",
			},
			wantRes: false,
		},
		{
			name: "no capital",
			args: args{
				input: "sawit@pr0",
			},
			wantRes: false,
		},
		{
			name: "no numeric",
			args: args{
				input: "sawit@pro",
			},
			wantRes: false,
		},
		{
			name: "no non alpha-numeric",
			args: args{
				input: "sawitpr0",
			},
			wantRes: false,
		},
		{
			name: "passed",
			args: args{
				input: "Sawit@Pr0",
			},
			wantRes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := validatePassword(tt.args.input)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("validatePassword() gotRes = %t, wantRes = %t", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_checkNewPhoneNumber(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		wantRes bool
		wantErr error
	}{
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
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344556").
					Return(repository.User{}, errors.New("expected GetUserByPhoneNumber error")).
					Times(1)
			},
			wantRes: false,
			wantErr: errors.New("expected GetUserByPhoneNumber error"),
		},
		{
			name: "not a new phone number",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344556").
					Return(repository.User{
						ID: 1,
					}, nil).
					Times(1)
			},
			wantRes: false,
			wantErr: nil,
		},
		{
			name: "a new phone number",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344556").
					Return(repository.User{}, nil).
					Times(1)
			},
			wantRes: true,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			gotRes, gotErr := checkNewPhoneNumber(tt.args.ctx, s, tt.args.phoneNumber)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("checkNewPhoneNumber() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("checkNewPhoneNumber() gotRes = %+v, wantRes = %+v", gotRes, tt.wantRes)
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}
