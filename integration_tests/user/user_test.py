import json
import pytest
import os
import http
import datetime
from pathlib2 import Path

from lib.database import database
from lib.api_request import api_request
from lib.json_utils import json_utils


def test_user():
    database.reinit()

    query_text = 'select * from system_user where uid = %s'
    second_query_text = 'select * from employee where uid = %s'

    requests = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'request.json'))['Requests']
    results = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'result.json'))['Results']

    response = api_request.ApiRequest.do_admin_request(requests[0])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    userUID = responseContent['Uid']
    query = database.query(query_text, userUID)
    assert len(query)==1
    employeeUID = query[0][1]
    assert query[0][2] == results[0]['Username']
    assert query[0][3] == results[0]['Password']
    query = database.query(second_query_text, employeeUID)
    assert query[0][2] == results[0]['WorkDocumentId']
    assert query[0][3] == results[0]['RoleId']

    response = api_request.ApiRequest.do_admin_request(requests[1])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    userUID = responseContent['Uid']
    query = database.query(query_text, userUID)
    assert len(query)==1
    employeeUID = query[0][1]
    assert query[0][2] == results[1]['Username']
    assert query[0][3] == results[1]['Password']
    query = database.query(second_query_text, employeeUID)
    assert query[0][2] == results[1]['WorkDocumentId']
    assert query[0][3] == results[1]['RoleId']

    req = requests[2]
    req['API'] = req['API'] + userUID
    response = api_request.ApiRequest.do_admin_request(req)
    assert response.status_code== 200
    query = database.query(query_text, userUID)
    assert len(query)==1
    employeeUID = query[0][1]
    assert query[0][2] == results[2]['Username']
    assert query[0][3] == results[2]['Password']
    query = database.query(second_query_text, employeeUID)
    assert query[0][2] == results[2]['WorkDocumentId']
    assert query[0][3] == results[2]['RoleId']

   
    response = api_request.ApiRequest.do_admin_request(requests[3])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    # Four ussers are added on initialization and two in this test
    assert len(responseContent['Users']) == 6

    req = requests[4]
    req['API'] = req['API'] + userUID
    response = api_request.ApiRequest.do_admin_request(req)
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert responseContent['Name'] == results[2]['Name']
    assert responseContent['Surname'] == results[2]['Surname'] 
    assert responseContent['RoleId'] == results[2]['RoleId']
  
    req = requests[5]
    req['API'] = req['API'] + userUID
    response = api_request.ApiRequest.do_admin_request(req)
    assert response.status_code== 200
    query = database.query(query_text, userUID)
    assert len(query)==0