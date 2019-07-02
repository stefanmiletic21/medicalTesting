package data

import (
	"context"
	"medicalTesting/db"
	"medicalTesting/dto"
)

var (
	CreateFilledTest = createFilledTest
	DeleteFilledTest = deleteFilledTest
	GetFilledTest    = getFilledTest
	GetFilledTests   = getFilledTests
)

func createFilledTest(ctx context.Context, createFilledTestRequest *dto.CreateFilledTestRequest) (uid string, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `insert into filled_test 
				(test_uid, examination_uid, answers) 
				values ($1, $2, $3)
				returning uid`

	rows, err := d.Query(ctx, query, createFilledTestRequest.TestUID, createFilledTestRequest.ExaminationUID, createFilledTestRequest.Answers)
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

func deleteFilledTest(ctx context.Context, filledTestUID string) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `delete from filled_test
				where uid = $1`

	_, err = d.Exec(ctx, query, filledTestUID)

	return
}

func getFilledTest(ctx context.Context, filledTestUID string) (test *dto.GetFilledTestResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `select 
				ft.test_uid as "TestUID",
				t.name as "TestName",
				ft.examination_uid as "ExaminationUID",
				ft.answers as "Answers"
				from filled_test ft
				join test t on (t.uid = ft.test_uid)
				where ft.uid = $1`

	rows, err := d.Query(ctx, query, filledTestUID)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		test = &dto.GetFilledTestResponse{}
		rr.ReadAllToStruct(test)
	}

	err = rr.Error()
	return
}

func getFilledTests(ctx context.Context) (filledTests *dto.GetFilledTestsResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)
	query := `select 
				ft.uid as "UID",
				e.examination_date as "ExaminationDate",
				e.patient_uid as "PatientUID",
				concat(p.name,' ', p.surname) as "PatientName",
				t.name as "TestName"
				from filled_test ft
				join examination e on (e.uid = ft.examination_uid)
				join test t on (t.uid = ft.test_uid)
				join patient pt on (pt.uid = e.patient_uid)
				join person p on (pt.person_uid = p.uid)`

	rows, err := d.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}
	filledTests = &dto.GetFilledTestsResponse{}
	filledTests.FilledTests = make([]*dto.FilledTestsInfo, 0, 0)
	for rr.ScanNext() {
		filledTestInfo := &dto.FilledTestsInfo{}
		rr.ReadAllToStruct(filledTestInfo)
		filledTests.FilledTests = append(filledTests.FilledTests, filledTestInfo)
	}
	err = rr.Error()
	return
}
