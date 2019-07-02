package data

import (
	"context"
	"medicalTesting/db"
	"medicalTesting/dto"
)

var (
	Login              = login
	CreateSession      = createSession
	GetSessionsForUser = getSessionsForUser
	RemoveSession      = removeSession
)

func login(ctx context.Context, name, passhash string) (session *dto.SessionInfo, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `
		select
			s.uid as UserUID, 
			e.role_id as Role, 
			s.username as Username
		from system_user s
		join employee e on (s.employee_uid = e.uid)
		where username = $1
		and password = $2`

	rows, err := d.Query(ctx, query, name, passhash)

	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		session = &dto.SessionInfo{}
		rr.ReadAllToStruct(session)
	}

	err = rr.Error()
	return
}

func createSession(ctx context.Context, userUID, token string) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `insert into login_session 
				(system_user_uid, token) 
				values ($1, $2)`

	_, err = d.Exec(ctx, query, userUID, token)
	return
}

func removeSession(ctx context.Context, token string) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `delete from login_session 
			  where token = $1`

	_, err = d.Exec(ctx, query, token)

	return
}

func getSessionsForUser(ctx context.Context, userUID string) (sessions []*dto.SessionInfo, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `
		select su.uid as UserUID, 
		e.role_id as Role, 
		su.username as Username,
		ls.token as Token
		from system_user su
		join employee e on (su.employee_uid = e.uid)
		join login_session ls on (su.uid = ls.system_user_uid)
		where ls.system_user_uid = $1`

	rows, err := d.Query(ctx, query, userUID)

	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}

	for rr.ScanNext() {
		session := &dto.SessionInfo{}
		rr.ReadAllToStruct(session)
		sessions = append(sessions, session)
	}

	err = rr.Error()
	return
}
