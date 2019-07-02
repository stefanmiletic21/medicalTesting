import json
import pytest
import os
import http
import datetime
from pathlib2 import Path

from lib.database import database
from lib.api_request import api_request
from lib.json_utils import json_utils


def test_test():
    database.reinit()

    query_text = 'select * from test where uid = %s'

    requests = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'request.json'))['Requests']
    results = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'result.json'))['Results']

    response = api_request.ApiRequest.do_doctor_request(requests[0])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    assert len(responseContent['Tests'])==1
    testUID = responseContent['Tests'][0]['Uid']

    req = requests[1]
    req['API'] = req['API'] + testUID
    response = api_request.ApiRequest.do_doctor_request(req)
    assert response.status_code== 200
    responseContent = json.loads(response.content)

    assert responseContent['Name'] == results[0]['Name']
    assert responseContent['Specialty'] == results[0]['Specialty']
    assert responseContent['Questions'] == results[0]['Questions']

    req = requests[2]
    req['API'] = req['API'] + testUID
    response = api_request.ApiRequest.do_doctor_request(req)
    assert response.status_code== 200
    query = database.query(query_text, testUID)
    assert len(query)==0
    

   