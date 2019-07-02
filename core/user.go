package core

import (
	"context"
	"encoding/json"
	"io"
	"medicalTesting/data"
	"medicalTesting/dto"
	"medicalTesting/logger"
	"medicalTesting/serverErr"
	"medicalTesting/utils"
)

var (
	CreateUser = createUser
	UpdateUser = updateUser
	RemoveUser = removeUser
	GetUser    = getUser
	GetUsers   = getUsers
)

func createUser(ctx context.Context, requestBody io.Reader) (response *dto.CreateUserResponse, err error) {
	request := &dto.CreateUserRequest{}
	response = &dto.CreateUserResponse{}
	err = json.NewDecoder(requestBody).Decode(request)
	if err != nil {
		logger.Warn("Request data could not be decoded: %v", err)
		err = serverErr.ErrBadRequest
		return
	}
	if request.EmployeeUID == "" {
		createPersonRequest := &dto.CreatePersonRequest{Address: request.Address, DateOfBirth: request.DateOfBirth, Email: request.Email, JMBG: request.JMBG, Name: request.Name, Surname: request.Surname}
		uid, err1 := data.CreatePerson(ctx, createPersonRequest)
		if err1 != nil {
			err = err1
			logger.Error("Couldn't create person: %v", err)
			return
		}
		personUID := uid
		createEmployeeRequest := &dto.CreateEmployeeRequest{PersonUID: personUID, RoleID: request.RoleID, WorkDocumentID: request.WorkDocumentID}
		uid, err1 = data.CreateEmployee(ctx, createEmployeeRequest)
		if err1 != nil {
			err = err1
			logger.Error("Couldn't create employee: %v", err)
			return
		}
		request.EmployeeUID = uid
	}
	request.Password = utils.GetPasswordHash(request.Password)
	uid, err := data.CreateUser(ctx, request)
	if err != nil {
		logger.Error("Couldn't create user: %v", err)
		err = serverErr.ErrInternal
		return
	}

	response.UID = uid
	return
}

func updateUser(ctx context.Context, userUID string, requestBody io.Reader) (err error) {
	request := &dto.UpdateUserRequest{}
	err = json.NewDecoder(requestBody).Decode(request)
	if err != nil {
		logger.Warn("Request data could not be decoded: %v", err)
		err = serverErr.ErrBadRequest
		return
	}

	updateEmployee := &dto.UpdateEmployeeRequest{WorkDocumentID: request.WorkDocumentID}
	employeeUID, err := data.UpdateEmployeeForUser(ctx, userUID, updateEmployee)
	if err != nil {
		logger.Error("Couldn't update employee for user with uid %v: %v", userUID, err)
		err = serverErr.ErrInternal
	}
	updatePerson := &dto.UpdatePersonRequest{Name: request.Name, Surname: request.Surname, JMBG: request.JMBG, Email: request.Email, Address: request.Address, DateOfBirth: request.DateOfBirth}
	err = data.UpdatePersonForEmployee(ctx, employeeUID, updatePerson)
	if err != nil {
		logger.Error("Couldn't update person for employee with uid %v: %v", employeeUID, err)
		err = serverErr.ErrInternal
	}

	request.Password = utils.GetPasswordHash(request.Password)
	err = data.UpdateUser(ctx, userUID, request)
	if err != nil {
		logger.Error("Couldn't update user with uid %v: %v", userUID, err)
		err = serverErr.ErrInternal
	}

	return
}

func removeUser(ctx context.Context, userUID string) (err error) {
	err = data.DeleteUser(ctx, userUID)
	if err != nil {
		logger.Error("Couldn't delete user with uid %v: %v", userUID, err)
		err = serverErr.ErrInternal
	}

	return
}

func getUser(ctx context.Context, userUID string) (response *dto.GetUserResponse, err error) {

	response, err = data.GetUser(ctx, userUID)
	if err != nil {
		logger.Error("Couldn't get user with uid %v: %v", userUID, err)
		err = serverErr.ErrInternal
		return
	}

	return
}

func getUsers(ctx context.Context) (response *dto.GetUsersResponse, err error) {
	response, err = data.GetUsers(ctx)
	if err != nil {
		logger.Error("Couldn't get users with uid %v: %v", err)
		err = serverErr.ErrInternal
	}
	return
}
