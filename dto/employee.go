package dto

import "time"

type CreateEmployeeResponse struct {
	UID string `json:"Uid"`
}

type CreateEmployeeRequest struct {
	PersonUID      string `json:"PersonUid"`
	Name           string
	Surname        string
	JMBG           string
	DateOfBirth    time.Time `json:",string"`
	Address        string
	Email          string
	WorkDocumentID string `json:"WorkDocumentId"`
	RoleID         int    `json:"RoleId"`
}

type UpdateEmployeeRequest struct {
	Name           string
	Surname        string
	JMBG           string
	DateOfBirth    time.Time `json:",string"`
	Address        string
	Email          string
	WorkDocumentID string `json:"WorkDocumentId"`
	RoleID         int    `json:"RoleId"`
}

type GetEmployeeResponse struct {
	UID            string `json:"Uid"`
	Name           string
	Surname        string
	JMBG           string
	DateOfBirth    time.Time `json:",string"`
	Address        string
	Email          string
	WorkDocumentID string `json:"WorkDocumentId"`
	RoleID         int    `json:"RoleId"`
}

type EmployeeBasicInfo struct {
	UID            string `json:"Uid"`
	Name           string
	Surname        string
	WorkDocumentID string `json:"WorkDocumentId"`
	RoleID         int    `json:"RoleId,omitempty"`
}

type GetEmployeesResponse struct {
	Employees []*EmployeeBasicInfo
}
