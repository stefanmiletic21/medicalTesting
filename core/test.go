package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"medicalTesting/data"
	"medicalTesting/dto"
	"medicalTesting/logger"
	"medicalTesting/serverErr"
	"medicalTesting/utils"
	"strconv"
)

var (
	CreateTest = createTest
	RemoveTest = removeTest
	GetTests   = getTests
	GetTest    = getTest
)

func createTest(ctx context.Context, request *http.Request) (err error) {
	name := request.FormValue("name")
	specialty, err := strconv.Atoi(request.FormValue("specialty"))
	if err != nil {
		logger.Warn("Couldn't get specialty from request")
		err = serverErr.ErrBadRequest
		return
	}
	file, header, err := request.FormFile("fileUpload")
	if err != nil {
		logger.Warn("Couldn't get file from request: %v", err)
		err = serverErr.ErrBadRequest
		return
	}

	defer file.Close()

	f, err := os.OpenFile("./"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)

	defer f.Close()
	defer os.Remove("./" + header.Filename)
	questions, err := utils.ParseFile(ctx, "./"+header.Filename)
	if err != nil {
		logger.Warn("Couldn't parse file %v", err)
		err = serverErr.ErrBadRequest
		return
	}

	marshaledQuestions, err := json.Marshal(questions)
	if err != nil {
		logger.Error("Couldn't marshall questions %v", err)
		return
	}
	err = data.CreateTest(ctx, name, specialty, marshaledQuestions)
	if err != nil {
		logger.Error("Couldn't create test %v", err)
		return
	}

	return
}

func removeTest(ctx context.Context, testUID string) (err error) {
	err = data.DeleteTest(ctx, testUID)
	if err != nil {
		logger.Error("Couldn't remove test with uid %v: %v", testUID, err)
		err = serverErr.ErrInternal
	}

	return
}

func getTest(ctx context.Context, testUID string) (response *dto.GetTestResponse, err error) {

	response, err = data.GetTest(ctx, testUID)
	if err != nil {
		logger.Error("Couldn't get test with uid %v: %v", testUID, err)
		err = serverErr.ErrInternal
		return
	}

	return
}

func getTests(ctx context.Context) (response *dto.GetTestsResponse, err error) {
	response, err = data.GetTests(ctx)
	if err != nil {
		logger.Error("Couldn't get tests: %v", err)
		err = serverErr.ErrInternal
		return
	}

	return
}
