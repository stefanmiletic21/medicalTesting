package data

import (
	"context"
	"medicalTesting/db"
	"medicalTesting/dto"
)

var (
	CreateExamination = createExamination
	DeleteExamination = deleteExamination
	GetExaminations   = getExaminations
	GetMyExaminations = getMyExaminations
)

func createExamination(ctx context.Context, request *dto.CreateExaminationRequest) (uid string, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `insert into examination 
				(doctor_uid, patient_uid, examination_date, doctor_full_name, patient_full_name) 
				values ($1, $2, $3,(
						SELECT CONCAT(p.name, ' ', p.surname)
						from employee e
						join person p on (e.person_uid = p.uid)
						where e.uid = $1
					),(
						SELECT CONCAT(p.name,' ', p.surname)
						from patient pa
						join person p on (pa.person_uid = p.uid)
						where pa.uid = $2
						)
				)
				returning uid`

	rows, err := d.Query(ctx, query, request.DoctorUID, request.PatientUID, request.ExaminationDate)
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

func deleteExamination(ctx context.Context, examinationUID string) (err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)

	query := `delete from examination
				where uid = $1`

	_, err = d.Exec(ctx, query, examinationUID)

	return
}

func getExaminations(ctx context.Context) (examinations *dto.GetExaminationsResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)
	query := `select 
				uid as "UID",
				doctor_uid as "DoctorUID",
				doctor_full_name as "DoctorFullName",
				patient_uid as "PatientUID",
				patient_full_name as "PatientFullName",
				examination_date as "ExaminationDate"
	 			from examination`

	rows, err := d.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}
	examinations = &dto.GetExaminationsResponse{}
	examinationInfos := make([]*dto.ExaminationInfo, 0)
	for rr.ScanNext() {
		examination := &dto.ExaminationInfo{}
		rr.ReadAllToStruct(examination)
		examinationInfos = append(examinationInfos, examination)
	}
	examinations.Examinations = examinationInfos
	err = rr.Error()
	return
}

func getMyExaminations(ctx context.Context, doctorUID string) (examinations *dto.GetExaminationsResponse, err error) {
	d := ctx.Value(db.RunnerKey).(db.Runner)
	query := `select 
				uid as "UID",
				doctor_uid as "DoctorUID",
				doctor_full_name as "DoctorFullName",
				patient_uid as "PatientUID",
				patient_full_name as "PatientFullName",
				examination_date as "ExaminationDate"
				from examination
				where doctor_uid = $1`

	rows, err := d.Query(ctx, query, doctorUID)
	if err != nil {
		return
	}
	defer rows.Close()

	rr, err := db.GetRowReader(rows)
	if err != nil {
		return
	}
	examinations = &dto.GetExaminationsResponse{}
	examinationInfos := make([]*dto.ExaminationInfo, 0)
	for rr.ScanNext() {
		examination := &dto.ExaminationInfo{}
		rr.ReadAllToStruct(examination)
		examinationInfos = append(examinationInfos, examination)
	}
	examinations.Examinations = examinationInfos
	err = rr.Error()
	return
}
