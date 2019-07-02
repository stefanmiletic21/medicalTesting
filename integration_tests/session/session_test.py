import json
import pytest
import os
import http
from pathlib2 import Path

from lib.database import database
from lib.api_request import api_request
from lib.json_utils import json_utils


def test_login_logout():
    database.reinit()

    query_text = 'select * from login_session  where system_user_uid = %s and token = %s'

    requests = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'request.json'))['Requests']

    response = api_request.ApiRequest.login_test_request(requests[0])
    assert response.status_code== 200
    responseContent = json.loads(response.content)
    userUID = responseContent['UserUID']
    auth = responseContent['Authorization']
    query = database.query(query_text, userUID, auth)
    assert len(query) == 1
    
    response = api_request.ApiRequest.login_test_request(requests[1],auth)
    assert response.status_code== 200
    query = database.query(query_text, userUID, auth)
    assert len(query) == 0
    response = api_request.ApiRequest.login_test_request(requests[2])
    assert response.status_code== 401