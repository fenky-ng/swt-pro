package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	errorHelper "github.com/fenky-ng/swt-pro/helper/error"
)

func Test_Repository_GetUserByID(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_GetUserByID] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		wantRes User
		wantErr error
	}{
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByID)).
					WithArgs(int64(1)).
					WillReturnError(errors.New("expected error"))
			},
			wantRes: User{},
			wantErr: errors.New("expected error"),
		},
		{
			name: "no data",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id", "phone_number", "password", "full_name"})

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByID)).
					WithArgs(int64(1)).
					WillReturnRows(resultRows)
			},
			wantRes: User{},
			wantErr: nil,
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id", "phone_number", "password", "full_name"}).
					AddRow(1, "+628223344556", "<password>", "Sawit")

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByID)).
					WithArgs(int64(1)).
					WillReturnRows(resultRows)
			},
			wantRes: User{
				ID:          1,
				PhoneNumber: "+628223344556",
				Password:    "<password>",
				FullName:    "Sawit",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			gotRes, gotErr := r.GetUserByID(tt.args.ctx, tt.args.userID)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("Repository.GetUserByID() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Repository.GetUserByID() gotRes = %+v, wantRes = %+v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_Repository_GetUserByPhoneNumber(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_GetUserByPhoneNumber] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
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
		wantRes User
		wantErr error
	}{
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByPhoneNumber)).
					WithArgs("+628223344556").
					WillReturnError(errors.New("expected error"))
			},
			wantRes: User{},
			wantErr: errors.New("expected error"),
		},
		{
			name: "no data",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id", "phone_number", "password", "full_name"})

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByPhoneNumber)).
					WithArgs("+628223344556").
					WillReturnRows(resultRows)
			},
			wantRes: User{},
			wantErr: nil,
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id", "phone_number", "password", "full_name"}).
					AddRow(1, "+628223344556", "<password>", "Sawit")

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByPhoneNumber)).
					WithArgs("+628223344556").
					WillReturnRows(resultRows)
			},
			wantRes: User{
				ID:          1,
				PhoneNumber: "+628223344556",
				Password:    "<password>",
				FullName:    "Sawit",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			gotRes, gotErr := r.GetUserByPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("Repository.GetUserByPhoneNumber() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Repository.GetUserByPhoneNumber() gotRes = %+v, wantRes = %+v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_Repository_IncreaseLoginCount(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_IncreaseLoginCount] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		wantErr error
	}{
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				sqlMock.ExpectExec(regexp.QuoteMeta(queryIncreaseLoginCount)).
					WithArgs(int64(1)).
					WillReturnError(errors.New("expected error"))
			},
			wantErr: errors.New("expected error"),
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				sqlMock.ExpectExec(regexp.QuoteMeta(queryIncreaseLoginCount)).
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			gotErr := r.IncreaseLoginCount(tt.args.ctx, tt.args.userID)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("Repository.IncreaseLoginCount() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
		})
	}
}

func Test_Repository_InsertUser(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_InsertUser] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx  context.Context
		data User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		wantRes int64
		wantErr error
	}{
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					PhoneNumber: "+628223344556",
					Password:    "<password>",
					FullName:    "Sawit",
				},
			},
			mock: func(fields *fields) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(queryInsertUser)).
					WithArgs("+628223344556", "<password>", "Sawit").
					WillReturnError(errors.New("expected error"))
			},
			wantRes: 0,
			wantErr: errors.New("expected error"),
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					PhoneNumber: "+628223344556",
					Password:    "<password>",
					FullName:    "Sawit",
				},
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id"}).
					AddRow(1)

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryInsertUser)).
					WithArgs("+628223344556", "<password>", "Sawit").
					WillReturnRows(resultRows)
			},
			wantRes: 1,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			gotRes, gotErr := r.InsertUser(tt.args.ctx, tt.args.data)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("Repository.InsertUser() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Repository.InsertUser() gotRes = %+v, wantRes = %+v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_Repository_UpdateUser(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_UpdateUser] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx  context.Context
		data User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		wantErr error
	}{
		{
			name: "no changes",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					ID:          1,
					PhoneNumber: "",
					FullName:    "",
				},
			},
			mock:    func(fields *fields) {},
			wantErr: nil,
		},
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					ID:          1,
					PhoneNumber: "+62812345678",
					FullName:    "New Name",
				},
			},
			mock: func(fields *fields) {
				sqlMock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(queryUpdateUser, "phone_number = $2, full_name = $3"))).
					WithArgs(int64(1), "+62812345678", "New Name").
					WillReturnError(errors.New("expected error"))
			},
			wantErr: errors.New("expected error"),
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					ID:          1,
					PhoneNumber: "+62812345678",
					FullName:    "New Name",
				},
			},
			mock: func(fields *fields) {
				sqlMock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(queryUpdateUser, "phone_number = $2, full_name = $3"))).
					WithArgs(int64(1), "+62812345678", "New Name").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			gotErr := r.UpdateUser(tt.args.ctx, tt.args.data)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("Repository.UpdateUser() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
		})
	}
}
