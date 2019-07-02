package data

import (
	"context"
	"medicalTesting/db"
	"medicalTesting/dto"
)

var (
	CreateTest = createTest
	DeleteTest = deleteTest
	GetTests   = getTests
	GetTest    = getTest
)

func createTest(ctx context.Context, name string, specialtyID int, questions []byte) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `insert into test 
				(name, specialty_id, questions) 
				values ($1, $2, $3)`

	_, err = d.Exec(ctx, query, name, specialtyID, questions)

	return
}

func deleteTest(ctx context.Context, testUID string) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `delete from test
				where uid = $1`

	_, err = d.Exec(ctx, query, testUID)

	return
}

func getTest(ctx context.Context, testUID string) (test *dto.GetTestResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `select 
				uid as "UID",
				name as "Name",
				specialty_id as "Specialty",
				questions->'Questions' as "Questions"
	 			from test 
				where uid = $1`

	rows, err := d.Query(ctx, query, testUID)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		test = &dto.GetTestResponse{}
		rr.ReadAllToStruct(test)
	}

	err = rr.Error()
	return
}

func getTests(ctx context.Context) (tests *dto.GetTestsResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)
	query := `select 
				uid as "UID",
				name as "Name",
				specialty_id as "Specialty"
	 			from test`

	rows, err := d.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}
	tests = &dto.GetTestsResponse{}
	testList := make([]*dto.TestInfo, 0)

	for rr.ScanNext() {
		test := &dto.TestInfo{}
		rr.ReadAllToStruct(test)
		testList = append(testList, test)
	}

	tests.Tests = testList
	err = rr.Error()
	return
}
