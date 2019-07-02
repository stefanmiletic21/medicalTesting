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
	CreateExamination = createExamination
	RemoveExamination = removeExamination
	GetExaminations   = getExaminations
	GetMyExaminations = getMyExaminations
)

func createExamination(ctx context.Context, requestBody io.Reader) (response *dto.CreateExaminationResponse, err error) {
	request := &dto.CreateExaminationRequest{}
	response = &dto.CreateExaminationResponse{}
	err = json.NewDecoder(requestBody).Decode(request)
	if err != nil {
		logger.Warn("Request data cannot be decoded: %v", err)
		err = serverErr.ErrBadRequest
		return
	}

	uid, err := data.CreateExamination(ctx, request)
	if err != nil {
		logger.Error("Couldn't create examination: %v", err)
		err = serverErr.ErrInternal
		return
	}

	response.UID = uid
	return
}

func removeExamination(ctx context.Context, examinationUID string) (err error) {
	err = data.DeleteExamination(ctx, examinationUID)
	if err != nil {
		logger.Error("Couldn't remove examination with uid %v: %v", examinationUID, err)
		err = serverErr.ErrInternal
	}

	return
}

func getExaminations(ctx context.Context) (response *dto.GetExaminationsResponse, err error) {
	response, err = data.GetExaminations(ctx)
	if err != nil {
		logger.Error("Couldn't get examinations: %v", err)
		err = serverErr.ErrInternal
		return
	}

	return
}

func getMyExaminations(ctx context.Context) (response *dto.GetExaminationsResponse, err error) {
	user := ctx.Value(UserKey).(*dto.SessionInfo)
	doctor, err := data.GetUser(ctx, user.UserUID)
	if err != nil {
		logger.Error("Couldn't get user with uid %v: %v", user.UserUID, err)
		err = serverErr.ErrInternal
		return
	}
	response, err = data.GetMyExaminations(ctx, doctor.UID)
	if err != nil {
		logger.Error("Couldn't get examinations for doctor with uid %v: %v", doctor.UID, err)
		err = serverErr.ErrInternal
		return
	}

	return
}
