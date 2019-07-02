package data

import (
	"context"
	"medicalTesting/db"
	"medicalTesting/dto"
)

var (
	CreatePerson            = createPerson
	UpdatePerson            = updatePerson
	UpdatePersonForEmployee = updatePersonForEmployee
	UpdatePersonForPatient  = updatePersonForPatient
	DeletePerson            = deletePerson
	GetPerson               = getPerson
	GetPersons              = getPersons
)

func createPerson(ctx context.Context, request *dto.CreatePersonRequest) (uid string, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `insert into person 
				(name, surname, JMBG, date_of_birth, address, email) 
				values ($1, $2, $3, $4, $5, $6)
				returning uid`

	rows, err := d.Query(ctx, query, request.Name, request.Surname, request.JMBG, request.DateOfBirth, request.Address, request.Email)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		uid = rr.ReadByIdxString(0)
	}

	err = rr.Error()
	return
}

func updatePerson(ctx context.Context, personUID string, request *dto.UpdatePersonRequest) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `update person set
				name = $1, surname = $2, JMBG = $3, date_of_birth = $4, address = $5, email = $6
				where uid = $7`

	_, err = d.Exec(ctx, query, request.Name, request.Surname, request.JMBG, request.DateOfBirth, request.Address, request.Email, personUID)

	return
}

func updatePersonForEmployee(ctx context.Context, employeeUID string, request *dto.UpdatePersonRequest) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `update person set
				name = $1, surname = $2, JMBG = $3, date_of_birth = $4, address = $5, email = $6
				where uid = (select person_uid from employee where uid = $7 limit 1)`

	_, err = d.Exec(ctx, query, request.Name, request.Surname, request.JMBG, request.DateOfBirth, request.Address, request.Email, employeeUID)

	return
}

func updatePersonForPatient(ctx context.Context, patientUID string, request *dto.UpdatePersonRequest) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `update person set
				name = $1, surname = $2, JMBG = $3, date_of_birth = $4, address = $5, email = $6
				where uid = (select person_uid from patient where uid = $7 limit 1)`

	_, err = d.Exec(ctx, query, request.Name, request.Surname, request.JMBG, request.DateOfBirth, request.Address, request.Email, patientUID)

	return
}

func deletePerson(ctx context.Context, personUID string) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `delete from person
				where uid = $1`

	_, err = d.Exec(ctx, query, personUID)

	return
}

func getPerson(ctx context.Context, personUID string) (person *dto.GetPersonResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `select 
				uid as "UID",
				name as "Name",
				surname as "Surname",
				jmbg as "JMBG",
				address as "Address",
				email as "Email",
				date_of_birth as "DateOfBirth"
	 			from person 
				where uid = $1`

	rows, err := d.Query(ctx, query, personUID)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		person = &dto.GetPersonResponse{}
		rr.ReadAllToStruct(person)
	}

	err = rr.Error()
	return
}

func getPersons(ctx context.Context) (persons *dto.GetPersonsResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)
	query := `select 
				uid as "UID",
				name as "Name",
				surname as "Surname",
				jmbg as "JMBG",
				address as "Address",
				email as "Email"
	 			from person`

	rows, err := d.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}
	persons = &dto.GetPersonsResponse{}
	for rr.ScanNext() {
		person := &dto.PersonBasicInfo{}
		rr.ReadAllToStruct(person)
		persons.Persons = append(persons.Persons, person)
	}
	err = rr.Error()
	return
}
