package handler

import (
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
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		wantErr error
	}{
		// {
		// 	// TODO fenky
		// },
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
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		wantErr error
	}{
		// {
		// 	// TODO fenky
		// },
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
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		wantErr error
	}{
		// {
		// 	// TODO fenky
		// },
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
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		wantErr error
	}{
		// {
		// 	// TODO fenky
		// },
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
			tt.fields.mockCtrl.Finish()
		})
	}
}
