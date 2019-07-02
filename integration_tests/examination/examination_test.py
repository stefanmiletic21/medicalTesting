import json
import pytest
import os
import http
import datetime
from pathlib2 import Path

from lib.database import database
from lib.api_request import api_request
from lib.json_utils import json_utils


def test_examination():
    database.reinit()

    query_text = 'select * from examination where uid = %s'

    requests = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'request.json'))['Requests']
    results = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'result.json'))['Results']

    response = api_request.ApiRequest.do_nurse_request(requests[0])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    examinationUID = responseContent['Uid']
    query = database.query(query_text, examinationUID)
    assert len(query)==1
    assert query[0][1] == results[0]['DoctorUid']
    assert query[0][2] == results[0]['DoctorName']
    assert query[0][3] == results[0]['PatientUid']
    assert query[0][4] == results[0]['PatientName']
    examinationDate = datetime.datetime.strptime(results[0]['ExaminationDate'], '%Y-%m-%dT%H:%M:%SZ')
    assert query[0][5] == examinationDate.date()

    response = api_request.ApiRequest.do_nurse_request(requests[1])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert len(responseContent['Examinations'])==2
    assert responseContent['Examinations'][1]['DoctorUid'] == results[0]['DoctorUid']
    assert responseContent['Examinations'][1]['DoctorFullName'] == results[0]['DoctorName']
    assert responseContent['Examinations'][1]['PatientUid'] == results[0]['PatientUid']
    assert responseContent['Examinations'][1]['PatientFullName'] == results[0]['PatientName']
    assert responseContent['Examinations'][1]['ExaminationDate'] == results[0]['ExaminationDate']
    
    response = api_request.ApiRequest.do_doctor_request(requests[2])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert len(responseContent['Examinations']) == 2

    req = requests[3]
    req['API'] = req['API'] + examinationUID
    response = api_request.ApiRequest.do_nurse_request(req)
    assert response.status_code== 200
    query = database.query(query_text, examinationUID)
    assert len(query)==0