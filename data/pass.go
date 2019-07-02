package data

import (
	"context"
	"medicalTesting/db"
)

var (
	CheckPassword  = checkPassword
	ChangePassword = changePassword
)

func checkPassword(ctx context.Context, pass, userUID string) (passed bool, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `select exists(select * from system_user where uid = $1 and password = $2)`

	rows, err := d.Query(ctx, query, userUID, pass)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		passed = true
	}

	err = rr.Error()
	return
}

func changePassword(ctx context.Context, newPass, userUID string) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `update system_user set
				password = $1
				where uid = $2`

	_, err = d.Exec(ctx, query, newPass, userUID)

	return
}
