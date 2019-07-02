import json
import pytest
import os
import http
import datetime
from pathlib2 import Path

from lib.database import database
from lib.api_request import api_request
from lib.json_utils import json_utils


def test_person():
    database.reinit()

    query_text = 'select * from person  where uid = %s'

    requests = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'request.json'))['Requests']
    results = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'result.json'))['Results']

    response = api_request.ApiRequest.do_admin_request(requests[0])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    personUID = responseContent['Uid']
    query = database.query(query_text, personUID)
    assert len(query)==1
    assert query[0][1] == results[0]['Name']
    assert query[0][2] == results[0]['Surname']
    assert query[0][3] == results[0]['JMBG']  
    date_of_birth = datetime.datetime.strptime(results[0]['DateOfBirth'], '%Y-%m-%dT%H:%M:%SZ')
    assert query[0][4] == date_of_birth.date()
    assert query[0][5] == results[0]['Address']
    assert query[0][6] == results[0]['Email']

    req  = requests[1]
    req['API'] = req['API'] + personUID
    response = api_request.ApiRequest.do_admin_request(req)
    assert response.status_code== 200
    query = database.query(query_text, personUID)
    assert len(query)==1
    assert query[0][1] == results[1]['Name']
    assert query[0][2] == results[1]['Surname']
    assert query[0][3] == results[1]['JMBG']  
    date_of_birth = datetime.datetime.strptime(results[1]['DateOfBirth'], '%Y-%m-%dT%H:%M:%SZ')
    assert query[0][4] == date_of_birth.date()
    assert query[0][5] == results[1]['Address']
    assert query[0][6] == results[1]['Email']

    response = api_request.ApiRequest.do_admin_request(requests[2])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    # Four persons are added on initialization and one in this test
    assert len(responseContent['Persons']) == 5
    assert responseContent['Persons'][4]['Name'] == results[1]['Name']
    assert responseContent['Persons'][4]['Surname'] == results[1]['Surname']
    assert responseContent['Persons'][4]['JMBG'] == results[1]['JMBG'] 

    req = requests[3]
    req['API'] = req['API'] + personUID
    response = api_request.ApiRequest.do_admin_request(req)
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert responseContent['Name'] == results[1]['Name']
    assert responseContent['Surname'] == results[1]['Surname']
    assert responseContent['JMBG'] == results[1]['JMBG'] 
    assert responseContent['DateOfBirth'] == results[1]['DateOfBirth'] 
    assert responseContent['Address'] == results[1]['Address'] 
    assert responseContent['Email'] == results[1]['Email'] 

    req = requests[4]
    req['API'] = req['API'] + personUID
    response = api_request.ApiRequest.do_admin_request(req)
    assert response.status_code== 200
    query = database.query(query_text, personUID)
    assert len(query)==0