package data

import (
	"context"
	"medicalTesting/db"
	"medicalTesting/dto"
	"medicalTesting/enum"
)

var (
	CreateEmployee        = createEmployee
	UpdateEmployee        = updateEmployee
	UpdateEmployeeForUser = updateEmployeeForUser
	DeleteEmployee        = deleteEmployee
	GetEmployee           = getEmployee
	GetEmployees          = getEmployees
	GetDoctors            = getDoctors
)

func createEmployee(ctx context.Context, request *dto.CreateEmployeeRequest) (uid string, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `insert into employee 
				(person_uid, work_document_id,role_id) 
				values ($1, $2,$3)
				returning uid`

	rows, err := d.Query(ctx, query, request.PersonUID, request.WorkDocumentID, request.RoleID)
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

func updateEmployee(ctx context.Context, employeeUID string, request *dto.UpdateEmployeeRequest) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `update employee set
				work_document_id = $1,
				role_id = $2
				where uid = $3`

	_, err = d.Exec(ctx, query, request.WorkDocumentID, request.RoleID, employeeUID)

	return
}

func deleteEmployee(ctx context.Context, employeeUID string) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `delete from employee
				where uid = $1`

	_, err = d.Exec(ctx, query, employeeUID)

	return
}

func getEmployee(ctx context.Context, employeeUID string) (employee *dto.GetEmployeeResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `select 
				e.uid as "UID",
				pe.name as "Name",
				pe.surname as "Surname",
				pe.jmbg as "JMBG",
				pe.address as "Address",
				pe.email as "Email",
				pe.date_of_birth as "DateOfBirth",
				e.work_document_id as "WorkDocumentID",
				e.role_id as "RoleID"
				from employee e
				join person pe on (e.person_uid = pe.uid) 
				where e.uid = $1`

	rows, err := d.Query(ctx, query, employeeUID)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		employee = &dto.GetEmployeeResponse{}
		rr.ReadAllToStruct(employee)
	}

	err = rr.Error()
	return
}

func getEmployees(ctx context.Context) (employees *dto.GetEmployeesResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)
	query := `select 
				e.uid as "UID",
				pe.name as "Name",
				pe.surname as "Surname",
				e.work_document_id as "WorkDocumentID",
				e.role_id as "RoleID"
				from employee e
				join person pe on (e.person_uid = pe.uid)`

	rows, err := d.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}
	employees = &dto.GetEmployeesResponse{}
	for rr.ScanNext() {
		employee := &dto.EmployeeBasicInfo{}
		rr.ReadAllToStruct(employee)
		employees.Employees = append(employees.Employees, employee)
	}
	err = rr.Error()
	return
}

func getDoctors(ctx context.Context) (employees *dto.GetEmployeesResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)
	query := `select 
				e.uid as "UID",
				pe.name as "Name",
				pe.surname as "Surname",
				e.work_document_id as "WorkDocumentID"
				from employee e
				join person pe on (e.person_uid = pe.uid)
				where role_id = $1`

	rows, err := d.Query(ctx, query, enum.RoleDoctor)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}
	employees = &dto.GetEmployeesResponse{}
	for rr.ScanNext() {
		employee := &dto.EmployeeBasicInfo{}
		rr.ReadAllToStruct(employee)
		employees.Employees = append(employees.Employees, employee)
	}
	err = rr.Error()
	return
}

func updateEmployeeForUser(ctx context.Context, userUID string, request *dto.UpdateEmployeeRequest) (uid string, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `update employee set
				work_document_id = $1
				where uid = (select employee_uid from system_user where uid = $2 limit 1)
				returning uid`

	rows, err := d.Query(ctx, query, request.WorkDocumentID, userUID)
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
