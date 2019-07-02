import json
import pytest
import os
import http
import datetime
from pathlib2 import Path

from lib.database import database
from lib.api_request import api_request
from lib.json_utils import json_utils


def test_filled_tests():
    database.reinit()

    query_text = 'select * from filled_test where uid = %s'

    requests = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'request.json'))['Requests']
    results = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'result.json'))['Results']

    response = api_request.ApiRequest.do_doctor_request(requests[0])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    filledTestUid = responseContent['Uid']
    query = database.query(query_text, filledTestUid)
    assert len(query)==1
    assert query[0][1] == results[0]['ExaminationUid']
    assert query[0][2] == results[0]['TestUid']
    assert query[0][3] == results[0]['Answers']

    response = api_request.ApiRequest.do_doctor_request(requests[1])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert len(responseContent['FilledTests']) ==  2
    assert responseContent['FilledTests'][1]['TestName'] == results[1]['TestName']
    assert responseContent['FilledTests'][1]['PatientName'] == results[1]['PatientName']
    assert responseContent['FilledTests'][1]['PatientUid'] == results[1]['PatientUid']

    req = requests[2]
    req['API'] = req['API'] + filledTestUid
    response = api_request.ApiRequest.do_doctor_request(req)
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert responseContent['TestName'] == results[2]['TestName']
    assert responseContent['TestUid'] == results[2]['TestUid']
    assert responseContent['ExaminationUID'] == results[2]['ExaminationUID']
    assert responseContent['Answers'] == results[2]['Answers']

    req = requests[3]
    req['API'] = req['API'] + filledTestUid
    response = api_request.ApiRequest.do_doctor_request(req)
    assert response.status_code== 200
    query = database.query(query_text, filledTestUid)
    assert len(query)==0
    