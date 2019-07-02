package dto

import (
	"time"
)

type CreatePatientResponse struct {
	UID string `json:"Uid"`
}

type CreatePatientRequest struct {
	PersonUID            string `json:"PersonUid"`
	Name                 string
	Surname              string
	JMBG                 string
	DateOfBirth          time.Time `json:",string"`
	Address              string
	Email                string
	MedicalRecordID      string    `json:"MedicalRecordId"`
	HealthCardID         string    `json:"HealthCardId"`
	HealthCardValidUntil time.Time `json:",string"`
}

type UpdatePatientRequest struct {
	Name                 string
	Surname              string
	JMBG                 string
	DateOfBirth          time.Time `json:",string"`
	Address              string
	Email                string
	MedicalRecordID      string    `json:"MedicalRecordId"`
	HealthCardID         string    `json:"HealthCardId"`
	HealthCardValidUntil time.Time `json:",string"`
}

type GetPatientResponse struct {
	UID                  string `json:"Uid"`
	Name                 string
	Surname              string
	JMBG                 string
	DateOfBirth          time.Time `json:",string"`
	Address              string
	Email                string
	MedicalRecordID      string    `json:"MedicalRecordId"`
	HealthCardID         string    `json:"HealthCardId"`
	HealthCardValidUntil time.Time `json:",string"`
}

type PatientBasicInfo struct {
	UID                  string `json:"Uid"`
	Name                 string
	Surname              string
	MedicalRecordID      string    `json:"MedicalRecordId"`
	HealthCardID         string    `json:"HealthCardId"`
	HealthCardValidUntil time.Time `json:",string"`
}

type GetPatientsResponse struct {
	Patients []*PatientBasicInfo
}
