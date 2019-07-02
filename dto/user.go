package dto

import "time"

type CreateUserResponse struct {
	UID string `json:"Uid"`
}

type CreateUserRequest struct {
	EmployeeUID    string `json:"EmployeeUid"`
	Name           string
	Surname        string
	JMBG           string
	DateOfBirth    time.Time `json:",string"`
	Address        string
	Email          string
	WorkDocumentID string `json:"WorkDocumentId"`
	RoleID         int    `json:"RoleId"`
	Username       string
	Password       string
}

type UpdateUserRequest struct {
	Name           string
	Surname        string
	JMBG           string
	DateOfBirth    time.Time `json:",string"`
	Address        string
	Email          string
	WorkDocumentID string `json:"WorkDocumentId"`
	RoleID         int    `json:"RoleId"`
	Username       string
	Password       string
}

type GetUserResponse struct {
	UID            string `json:"Uid"`
	Name           string
	Surname        string
	JMBG           string
	DateOfBirth    time.Time `json:",string"`
	Address        string
	Email          string
	WorkDocumentID string `json:"WorkDocumentId"`
	RoleID         int    `json:"RoleId"`
	Username       string
}

type UserBasicInfo struct {
	UID            string `json:"Uid"`
	Name           string
	Surname        string
	WorkDocumentID string `json:"WorkDocumentId"`
	RoleID         int    `json:"RoleId"`
	Username       string
}

type GetUsersResponse struct {
	Users []*UserBasicInfo
}
