package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fenky-ng/swt-pro/constant"
	"github.com/fenky-ng/swt-pro/generated"
	"github.com/fenky-ng/swt-pro/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Register
// (POST /register)
func (s *Server) Register(ctx echo.Context) error {
	var (
		funcName = "Register"
		request  generated.RegistrationRequest
		response generated.RegistrationResponse
	)

	// decode request body
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		log.Errorf("[%s] Decode error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeUnmarshal, []string{"Bad request"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// validate registration request
	requestValidationErrors := validateRegistration(request)
	if len(requestValidationErrors) != 0 {
		response.Header = generateResponseHeader(constant.ErrorCodeValidation, requestValidationErrors, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// check whether phone number is already registered or not
	isNewPhoneNumber, err := checkNewPhoneNumber(ctx.Request().Context(), s, request.PhoneNumber)
	if err != nil {
		log.Errorf("[%s] checkNewPhoneNumber error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeDatabase, []string{"System error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if !isNewPhoneNumber {
		response.Header = generateResponseHeader(constant.ErrorCodeValidation, []string{"Phone number is already registered"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// hash and salt the password
	salt, err := hashAndSalt(request.Password)
	if err != nil {
		log.Errorf("[%s] hashAndSalt error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeHashAndSalt, []string{"There was an error when handling password"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// insert user to db
	id, err := s.Repository.InsertUser(ctx.Request().Context(), repository.User{
		PhoneNumber: request.PhoneNumber,
		Password:    salt,
		FullName:    request.FullName,
	})
	if err != nil {
		log.Errorf("[%s] InsertUser error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeDatabase, []string{"System error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header = generateResponseHeader(0, nil, true)
	response.Data = &generated.RegistrationResponseData{
		Id: id,
	}

	return ctx.JSON(http.StatusOK, response)
}

// Login
// (POST /login)
func (s *Server) Login(ctx echo.Context) error {
	var (
		funcName = "Login"
		request  generated.LoginRequest
		response generated.LoginResponse
	)

	// decode request body
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		log.Errorf("[%s] Decode error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeUnmarshal, []string{"Bad request"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// get user from db by phone number
	user, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), request.PhoneNumber)
	if err != nil {
		log.Errorf("[%s] GetUserByPhoneNumber error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeDatabase, []string{"System error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// check whether user exists or not
	if user.ID == 0 {
		response.Header = generateResponseHeader(constant.ErrorCodeValidation, []string{"Phone number is not registered"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// check password
	if !comparePasswords(user.Password, request.Password) {
		response.Header = generateResponseHeader(constant.ErrorCodeValidation, []string{"Wrong password"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// generate jwt token
	jwtToken, err := generateJwtToken(user)
	if err != nil {
		log.Errorf("[%s] generateJwtToken error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeJWT, []string{"System error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// increase login count
	err = s.Repository.IncreaseLoginCount(ctx.Request().Context(), user.ID)
	if err != nil {
		log.Errorf("[%s] IncreaseLoginCount error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeDatabase, []string{"System error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header = generateResponseHeader(0, nil, true)
	response.Data = &generated.LoginResponseData{
		Id:  user.ID,
		Jwt: jwtToken,
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetProfile
// (GET /profile)
func (s *Server) GetProfile(ctx echo.Context) error {
	var (
		funcName = "GetProfile"
		response generated.GetProfileResponse
	)

	// get session claims
	sessionClaims, err := getSessionClaims(ctx)
	if err != nil {
		response.Header = generateResponseHeader(constant.ErrorCodeAuthorization, []string{err.Error()}, false)
		return ctx.JSON(http.StatusForbidden, response)
	}

	// get user from db by id
	user, err := s.Repository.GetUserByID(ctx.Request().Context(), sessionClaims.UserID)
	if err != nil {
		log.Errorf("[%s] GetUserByID error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeDatabase, []string{"System error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header = generateResponseHeader(0, nil, true)
	response.Data = &generated.GetProfileResponseData{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, response)
}

// UpdateProfile
// (PATCH /profile)
func (s *Server) UpdateProfile(ctx echo.Context) error {
	var (
		funcName = "UpdateProfile"
		request  generated.UpdateProfileRequest
		response generated.UpdateProfileResponse
	)

	// get session claims
	sessionClaims, err := getSessionClaims(ctx)
	if err != nil {
		response.Header = generateResponseHeader(constant.ErrorCodeAuthorization, []string{err.Error()}, false)
		return ctx.JSON(http.StatusForbidden, response)
	}

	// decode request body
	err = json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		log.Errorf("[%s] Decode error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeUnmarshal, []string{"Bad request"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	var (
		changesCount int
		phoneNumber  string
		fullName     string
	)

	if request.PhoneNumber != nil && *request.PhoneNumber != "" {
		phoneNumber = *request.PhoneNumber

		// validate phone number
		if errorMessages := validatePhoneNumber(phoneNumber); len(errorMessages) != 0 {
			response.Header = generateResponseHeader(constant.ErrorCodeValidation, errorMessages, false)
			return ctx.JSON(http.StatusBadRequest, response)
		}

		// check whether phone number is already registered or not
		user, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), phoneNumber)
		if err != nil {
			log.Errorf("[%s] GetUserByPhoneNumber error: %s", funcName, err.Error())
			response.Header = generateResponseHeader(constant.ErrorCodeDatabase, []string{"System error"}, false)
			return ctx.JSON(http.StatusInternalServerError, response)
		}
		if user.ID != 0 && user.ID != sessionClaims.UserID {
			response.Header = generateResponseHeader(constant.ErrorCodeValidation, []string{"Phone number is already registered"}, false)
			return ctx.JSON(http.StatusConflict, response)
		}

		changesCount++
	}
	if request.FullName != nil && *request.FullName != "" {
		fullName = *request.FullName

		// validate full name
		if errorMessages := validateFullName(fullName); len(errorMessages) != 0 {
			response.Header = generateResponseHeader(constant.ErrorCodeValidation, errorMessages, false)
			return ctx.JSON(http.StatusBadRequest, response)
		}

		changesCount++
	}
	if changesCount == 0 {
		response.Header = generateResponseHeader(constant.ErrorCodeValidation, []string{"No changes"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	err = s.Repository.UpdateUser(ctx.Request().Context(), repository.User{
		ID:          sessionClaims.UserID,
		PhoneNumber: phoneNumber,
		FullName:    fullName,
	})
	if err != nil {
		log.Errorf("[%s] UpdateUser error: %s", funcName, err.Error())
		response.Header = generateResponseHeader(constant.ErrorCodeDatabase, []string{"System error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header = generateResponseHeader(0, nil, true)

	return ctx.JSON(http.StatusOK, response)
}
