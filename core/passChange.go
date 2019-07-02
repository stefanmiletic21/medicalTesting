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
	ChangePass = changePass
)

func changePass(ctx context.Context, requestBody io.Reader, userUID string) (err error) {
	request := &dto.ChangePassRequest{}
	err = json.NewDecoder(requestBody).Decode(request)
	if err != nil {
		logger.Warn("Request data could not be decoded: %v")
		err = serverErr.ErrBadRequest
		return
	}

	passed, err := data.CheckPassword(ctx, utils.GetPasswordHash(request.OldPass), userUID)
	if err != nil {
		logger.Error("Couldn't check old password: %v")
		err = serverErr.ErrInternal
		return
	}
	if !passed {
		logger.Warn("User with uid %v is trying to change pass with bad authentication", userUID)
		err = serverErr.ErrNotAuthenticated
		return
	}
	err = data.ChangePassword(ctx, utils.GetPasswordHash(request.NewPass), userUID)
	if err != nil {
		logger.Error("Couldn't change password: %v", err)
		err = serverErr.ErrInternal
		return
	}
	return
}
