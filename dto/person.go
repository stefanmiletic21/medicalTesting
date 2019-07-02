package dto

import (
	"time"
)

type CreatePersonResponse struct {
	UID string `json:"Uid"`
}

type CreatePersonRequest struct {
	Name        string
	Surname     string
	JMBG        string
	DateOfBirth time.Time `json:",string"`
	Address     string
	Email       string
}

type UpdatePersonRequest struct {
	Name        string
	Surname     string
	JMBG        string
	DateOfBirth time.Time `json:",string"`
	Address     string
	Email       string
}

type GetPersonResponse struct {
	UID         string `json:"Uid"`
	Name        string
	Surname     string
	JMBG        string
	DateOfBirth time.Time `json:",string"`
	Address     string
	Email       string
}

type PersonBasicInfo struct {
	UID     string `json:"Uid"`
	Name    string
	Surname string
	JMBG    string
}

type GetPersonsResponse struct {
	Persons []*PersonBasicInfo
}
