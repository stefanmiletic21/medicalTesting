package data

import (
	"context"
	"medicalTesting/db"
	"medicalTesting/dto"
)

var (
	CreateUser = createUser
	UpdateUser = updateUser
	DeleteUser = deleteUser
	GetUser    = getUser
	GetUsers   = getUsers
)

func createUser(ctx context.Context, request *dto.CreateUserRequest) (uid string, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `insert into system_user 
				(employee_uid, username,password) 
				values ($1, $2, $3)
				returning uid`

	rows, err := d.Query(ctx, query, request.EmployeeUID, request.Username, request.Password)
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

func updateUser(ctx context.Context, userUID string, request *dto.UpdateUserRequest) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `update system_user set
				username = $1,
				password = $2
				where uid = $3`

	_, err = d.Exec(ctx, query, request.Username, request.Password, userUID)

	return
}

func deleteUser(ctx context.Context, userUID string) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `delete from system_user where uid = $1`

	_, err = d.Exec(ctx, query, userUID)

	return
}

func getUser(ctx context.Context, userUID string) (user *dto.GetUserResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `select 
				u.uid as "UID",
				pe.name as "Name",
				pe.surname as "Surname",
				pe.jmbg as "JMBG",
				pe.address as "Address",
				pe.email as "Email",
				pe.date_of_birth as "DateOfBirth",
				e.work_document_id as "WorkDocumentID",
				e.role_id as "RoleID",
				u.username as "Username",
				u.password as "Password"
				from system_user u
				join employee e on (u.employee_uid = e.uid)
				join person pe on (e.person_uid = pe.uid) 
				where u.uid = $1`

	rows, err := d.Query(ctx, query, userUID)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		user = &dto.GetUserResponse{}
		rr.ReadAllToStruct(user)
	}

	err = rr.Error()
	return
}

func getUsers(ctx context.Context) (users *dto.GetUsersResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)
	query := `select 
				u.uid as "UID",
				pe.name as "Name",
				pe.surname as "Surname",
				e.work_document_id as "WorkDocumentID",
				e.role_id as "RoleID",
				u.username as "Username",
				u.password as "Password"
				from system_user u
				join employee e on (u.employee_uid = e.uid)
				join person pe on (e.person_uid = pe.uid)
				order by u.uid asc`

	rows, err := d.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}
	users = &dto.GetUsersResponse{}
	for rr.ScanNext() {
		user := &dto.UserBasicInfo{}
		rr.ReadAllToStruct(user)
		users.Users = append(users.Users, user)
	}
	err = rr.Error()
	return
}
