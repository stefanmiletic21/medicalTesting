package dto

import (
	"encoding/json"
	"medicalTesting/enum"
)

type Question struct {
	Question string
	Type     enum.QuestionType
	Answers  []string
}

type TestQuestions struct {
	Questions []*Question
}

type GetTestsResponse struct {
	Tests []*TestInfo
}

type TestInfo struct {
	UID       string `json:"Uid"`
	Name      string
	Specialty int
}

type GetTestResponse struct {
	UID       string `json:"Uid"`
	Name      string
	Specialty int
	Questions json.RawMessage
}
