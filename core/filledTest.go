package core

import (
	"context"
	"encoding/json"
	"io"
	"medicalTesting/data"
	"medicalTesting/dto"
	"medicalTesting/logger"
	"medicalTesting/serverErr"
)

var (
	CreateFilledTest = createFilledTest
	RemoveFilledTest = removeFilledTest
	GetFilledTest    = getFilledTest
	GetFilledTests   = getFilledTests
)

func createFilledTest(ctx context.Context, requestBody io.Reader) (response *dto.CreateFilledTestResponse, err error) {
	request := &dto.CreateFilledTestRequest{}
	response = &dto.CreateFilledTestResponse{}

	err = json.NewDecoder(requestBody).Decode(request)
	if err != nil {
		logger.Warn("Request data cannot be decoded: %v", err)
		err = serverErr.ErrBadRequest
		return
	}

	uid, err := data.CreateFilledTest(ctx, request)
	if err != nil {
		logger.Error("Couldn't create filled test: %v", err)
		err = serverErr.ErrInternal
		return
	}

	response.UID = uid
	return
}

func removeFilledTest(ctx context.Context, filledTestUID string) (err error) {
	err = data.DeleteFilledTest(ctx, filledTestUID)
	if err != nil {
		logger.Error("Couldn't remove filled test with uid %v: %v", filledTestUID, err)
		err = serverErr.ErrInternal
	}

	return
}

func getFilledTest(ctx context.Context, filledTestUID string) (response *dto.GetFilledTestResponse, err error) {

	response, err = data.GetFilledTest(ctx, filledTestUID)
	if err != nil {
		logger.Error("Couldn't get filled test with uid %v: %v", filledTestUID, err)
		err = serverErr.ErrInternal
		return
	}

	return
}

func getFilledTests(ctx context.Context) (response *dto.GetFilledTestsResponse, err error) {
	response, err = data.GetFilledTests(ctx)
	if err != nil {
		logger.Error("Couldn't remove filled tests: %v", err)
		err = serverErr.ErrInternal
		return
	}

	return
}
