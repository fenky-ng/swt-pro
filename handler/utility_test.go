package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/fenky-ng/swt-pro/constant"
	"github.com/fenky-ng/swt-pro/generated"
	errorHelper "github.com/fenky-ng/swt-pro/helper/error"
	"github.com/fenky-ng/swt-pro/model"
	"github.com/fenky-ng/swt-pro/repository"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func Test_hashAndSalt(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr error
	}{
		{
			name: "password is too long",
			args: args{
				input: "1234567890" +
					"1234567890" +
					"1234567890" +
					"1234567890" +
					"1234567890" +
					"1234567890" +
					"1234567890" +
					"123",
			},
			wantRes: "",
			wantErr: errors.New("bcrypt: password length exceeds 72 bytes"),
		},
		{
			name: "passed",
			args: args{
				input: "Sawit@123",
			},
			wantRes: "let's say this is a hashed token",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, gotErr := hashAndSalt(tt.args.input)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("hashAndSalt() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if (len(gotRes) != 0) != (len(tt.wantRes) != 0) {
				t.Errorf("hashAndSalt() gotRes = %s, wantRes = %s", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_comparePasswords(t *testing.T) {
	type args struct {
		hashedPassword string
		plainPassword  string
	}
	tests := []struct {
		name    string
		args    args
		wantRes bool
	}{
		{
			name: "not match",
			args: args{
				hashedPassword: "$2a$04$sBtgSQFodfK4P/zPEV/uPuU7O3zwE6rQ8HyIFWRqKIBjwe5E1.CyC",
				plainPassword:  "Sawit@Pro",
			},
			wantRes: false,
		},
		{
			name: "match",
			args: args{
				hashedPassword: "$2a$04$sBtgSQFodfK4P/zPEV/uPuU7O3zwE6rQ8HyIFWRqKIBjwe5E1.CyC",
				plainPassword:  "Sawit@123",
			},
			wantRes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := comparePasswords(tt.args.hashedPassword, tt.args.plainPassword)
			if gotRes != tt.wantRes {
				t.Errorf("comparePasswords() gotRes = %t, wantRes = %t", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_generateJwtToken(t *testing.T) {
	type args struct {
		user repository.User
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr error
	}{
		{
			name: "passed",
			args: args{
				user: repository.User{
					ID:          1,
					PhoneNumber: "+628223344556",
				},
			},
			wantRes: "let's say this is a jwt token",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, gotErr := generateJwtToken(tt.args.user)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("generateJwtToken() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			if (len(gotRes) != 0) != (len(tt.wantRes) != 0) {
				t.Errorf("generateJwtToken() gotRes = %s, wantRes = %s", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_getSessionClaims(t *testing.T) {
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantRes model.SessionClaims
		wantErr error
	}{
		{
			name: "no jwt token",
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					res := httptest.NewRecorder()
					return echo.New().NewContext(req, res)
				}(),
			},
			wantRes: model.SessionClaims{},
			wantErr: errors.New("JWT token not found"),
		},
		{
			name: "token is expired",
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE1NDEwNDAsImlzcyI6InN3dC1wcm8iLCJ1c2VyX2lkIjoxLCJwaG9uZV9udW1iZXIiOiIrNjI4MjIzMzQ0NTUxIn0.eDQAImZTyrX9nNbAwIV_9FMT74EbBCnxMmOQ9fS1U7t7_Alx8f718TTzj6DK84g6CCaafRG1d576RApjBb0UA8t0WRmi6iMcvozN5cZ3dAz5aaHZPaJan7wPB1HhxF1lBEZZAds-EqHQghgtyrzlMlYQ92Lz-etkpRMevxuMGSxAsj9Pnp5azC_rmUg9XviOvxp35P0azm2WeD6c6h4BEeuaFGrwKsvH55OFyAg8VZou8NrIQKZsyIxKjjDMUsMQLcbcmeYaH2583HyjI9jMH0GzWc5fB-U4qCiwq2k9nsdfj2i2LbytN9CzxsdDIfHbrR74mcpyKjHxHJzSzGieNB8ddt-rbrtpYfPYTIXNG-duTQ8iK9p4-xSfrJ-RiIURsB_l5ipjmKgCFwiCO49SA5JsFoyfH5CirRyM4Psf4y5ayK0RO8n5B5D55QwQHf4LjTUBy5bW7pheHC4ZWz7S2YGIDZCr76iMuBSqYg0LHI7oT7M_aVVud7gNXbRDV2k1Pd8B30GStUThRQJd3-xuLkMpchCTRAZVZjYAeJ4Z6iaQBX3yz8sT38aZCH8u5D3PencOFS-VfHYsAUGBXHk-xZWF13QChUwtf-Yiab9l8CWSK1WtxbsZ5PdGha1aCtOgRk6iAIze-4o91kriYd-Bv8WG-wIBXGl92IHBrLUa-yY")
					res := httptest.NewRecorder()
					return echo.New().NewContext(req, res)
				}(),
			},
			wantRes: model.SessionClaims{},
			wantErr: errors.New("Session is expired"),
		},
		{
			name: "passed",
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					jwt, _ := generateJwtToken(repository.User{
						ID:          1,
						PhoneNumber: "+628223344556",
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					return echo.New().NewContext(req, res)
				}(),
			},
			wantRes: model.SessionClaims{
				StandardClaims: jwt.StandardClaims{
					Issuer: constant.ApplicationName,
				},
				UserID:      1,
				PhoneNumber: "+628223344556",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, gotErr := getSessionClaims(tt.args.ctx)
			if errorHelper.GetErrorMessage(gotErr) != errorHelper.GetErrorMessage(tt.wantErr) {
				t.Errorf("getSessionClaims() gotErr = %s, wantErr = %s", errorHelper.GetErrorMessage(gotErr), errorHelper.GetErrorMessage(tt.wantErr))
			}
			gotRes.ExpiresAt = 0
			if gotRes != tt.wantRes {
				t.Errorf("getSessionClaims() gotRes = %+v, wantRes = %+v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_generateResponseHeader(t *testing.T) {
	type args struct {
		errorCode     int
		errorMessages []string
		successful    bool
	}
	tests := []struct {
		name    string
		args    args
		wantRes generated.ResponseHeader
	}{
		{
			name: "error code",
			args: args{
				errorCode: constant.ErrorCodeGeneral,
			},
			wantRes: generated.ResponseHeader{
				ErrorCode: func() *int {
					res := constant.ErrorCodeGeneral
					return &res
				}(),
			},
		},
		{
			name: "error message",
			args: args{
				errorMessages: []string{"expected error"},
			},
			wantRes: generated.ResponseHeader{
				ErrorMessages: func() *[]string {
					res := []string{"expected error"}
					return &res
				}(),
			},
		},
		{
			name: "successful",
			args: args{
				successful: true,
			},
			wantRes: generated.ResponseHeader{
				Successful: func() *bool {
					res := true
					return &res
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := generateResponseHeader(tt.args.errorCode, tt.args.errorMessages, tt.args.successful)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("generateResponseHeader() gotRes = %+v, wantRes = %+v", gotRes, tt.wantRes)
			}
		})
	}
}
