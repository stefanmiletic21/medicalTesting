import json
import pytest
import os
import http
import datetime
from pathlib2 import Path

from lib.database import database
from lib.api_request import api_request
from lib.json_utils import json_utils


def test_patient():
    database.reinit()

    query_text = 'select * from patient where uid = %s'
    second_query_text = 'select * from person where uid = %s'

    requests = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'request.json'))['Requests']
    results = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'result.json'))['Results']

    response = api_request.ApiRequest.do_nurse_request(requests[0])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    employeeUID = responseContent['Uid']
    query = database.query(query_text, employeeUID)
    assert len(query)==1
    personUID = query[0][1]
    assert query[0][2] == results[0]['MedicalRecordId']
    assert query[0][3] == results[0]['HealthCardId']  
    query = database.query(second_query_text, personUID)
    assert query[0][1] == results[0]['Name']
    assert query[0][2] == results[0]['Surname']
    assert query[0][3] == results[0]['JMBG']  

    response = api_request.ApiRequest.do_nurse_request(requests[1])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    patientUID = responseContent['Uid']
    query = database.query(query_text, patientUID)
    assert len(query)==1
    personUID = query[0][1]
    assert query[0][2] == results[1]['MedicalRecordId']
    assert query[0][3] == results[1]['HealthCardId']  
    query = database.query(second_query_text, personUID)
    assert query[0][1] == results[1]['Name']
    assert query[0][2] == results[1]['Surname']
    assert query[0][3] == results[1]['JMBG']

    req = requests[2]
    req['API'] = req['API'] + patientUID
    response = api_request.ApiRequest.do_nurse_request(req)
    assert response.status_code== 200
    query = database.query(query_text, patientUID)
    assert len(query)==1
    personUID = query[0][1]
    assert query[0][2] == results[2]['MedicalRecordId']
    assert query[0][3] == results[2]['HealthCardId']  
    query = database.query(second_query_text, personUID)
    assert query[0][1] == results[2]['Name']
    assert query[0][2] == results[2]['Surname']
    assert query[0][3] == results[2]['JMBG'] 
   
    response = api_request.ApiRequest.do_nurse_request(requests[3])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    #Two patients are added on initialization and two now
    assert len(responseContent['Patients']) == 4
    assert responseContent['Patients'][3]['MedicalRecordId'] == results[2]['MedicalRecordId']
    assert responseContent['Patients'][3]['Name'] == results[2]['Name']
    assert responseContent['Patients'][3]['Surname'] == results[2]['Surname'] 
    assert responseContent['Patients'][3]['HealthCardId'] == results[2]['HealthCardId']

    req = requests[4]
    req['API'] = req['API'] + patientUID
    response = api_request.ApiRequest.do_nurse_request(req)
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert responseContent['Name'] == results[2]['Name']
    assert responseContent['Surname'] == results[2]['Surname'] 
    assert responseContent['HealthCardId'] == results[2]['HealthCardId']
    assert responseContent['JMBG'] == results[2]['JMBG']
    assert responseContent['Email'] == results[2]['Email']
  
    req = requests[5]
    req['API'] = req['API'] + patientUID
    response = api_request.ApiRequest.do_nurse_request(req)
    assert response.status_code== 200
    query = database.query(query_text, patientUID)
    assert len(query)==0