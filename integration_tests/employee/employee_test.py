import json
import pytest
import os
import http
import datetime
from pathlib2 import Path

from lib.database import database
from lib.api_request import api_request
from lib.json_utils import json_utils


def test_employee():
    database.reinit()

    query_text = 'select * from employee where uid = %s'
    second_query_text = 'select * from person where uid = %s'

    requests = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'request.json'))['Requests']
    results = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'result.json'))['Results']

    response = api_request.ApiRequest.do_admin_request(requests[0])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    employeeUID = responseContent['Uid']
    query = database.query(query_text, employeeUID)
    assert len(query)==1
    personUID = query[0][1]
    assert query[0][2] == results[0]['WorkDocumentId']
    assert query[0][3] == results[0]['RoleId']  
    query = database.query(second_query_text, personUID)
    assert query[0][1] == results[0]['Name']
    assert query[0][2] == results[0]['Surname']
    assert query[0][3] == results[0]['JMBG']  

    response = api_request.ApiRequest.do_admin_request(requests[1])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    employeeUID = responseContent['Uid']
    query = database.query(query_text, employeeUID)
    assert len(query)==1
    personUID = query[0][1]
    assert query[0][2] == results[1]['WorkDocumentId']
    assert query[0][3] == results[1]['RoleId']  
    query = database.query(second_query_text, personUID)
    assert query[0][1] == results[1]['Name']
    assert query[0][2] == results[1]['Surname']
    assert query[0][3] == results[1]['JMBG']

    req = requests[2]
    req['API'] = req['API'] + employeeUID
    response = api_request.ApiRequest.do_admin_request(req)
    assert response.status_code== 200
    query = database.query(query_text, employeeUID)
    assert len(query)==1
    personUID = query[0][1]
    assert query[0][2] == results[2]['WorkDocumentId']
    assert query[0][3] == results[2]['RoleId']  
    query = database.query(second_query_text, personUID)
    assert query[0][1] == results[2]['Name']
    assert query[0][2] == results[2]['Surname']
    assert query[0][3] == results[2]['JMBG'] 
   
    response = api_request.ApiRequest.do_admin_request(requests[3])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    # Four employees are added on initialization and two in this test
    assert len(responseContent['Employees']) == 6
    assert responseContent['Employees'][5]['WorkDocumentId'] == results[2]['WorkDocumentId']
    assert responseContent['Employees'][5]['Name'] == results[2]['Name']
    assert responseContent['Employees'][5]['Surname'] == results[2]['Surname'] 
    assert responseContent['Employees'][5]['RoleId'] == results[2]['RoleId']

    req = requests[4]
    req['API'] = req['API'] + employeeUID
    response = api_request.ApiRequest.do_admin_request(req)
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert responseContent['Name'] == results[2]['Name']
    assert responseContent['Surname'] == results[2]['Surname'] 
    assert responseContent['RoleId'] == results[2]['RoleId']
    assert responseContent['JMBG'] == results[2]['JMBG']
    assert responseContent['Email'] == results[2]['Email']

    response = api_request.ApiRequest.do_nurse_request(requests[5])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert len(responseContent['Employees']) == 2
    assert responseContent['Employees'][0]['WorkDocumentId'] == results[3]['Employees'][0]['WorkDocumentId']
    assert responseContent['Employees'][0]['Name'] == results[3]['Employees'][0]['Name']
    assert responseContent['Employees'][0]['Surname'] == results[3]['Employees'][0]['Surname'] 
    assert responseContent['Employees'][1]['WorkDocumentId'] == results[3]['Employees'][1]['WorkDocumentId']
    assert responseContent['Employees'][1]['Name'] == results[3]['Employees'][1]['Name']
    assert responseContent['Employees'][1]['Surname'] == results[3]['Employees'][1]['Surname'] 

    req = requests[6]
    req['API'] = req['API'] + employeeUID
    response = api_request.ApiRequest.do_admin_request(req)
    assert response.status_code== 200
    query = database.query(query_text, employeeUID)
    assert len(query)==0