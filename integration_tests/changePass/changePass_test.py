import json
import pytest
import os
import http
from pathlib2 import Path

from lib.database import database
from lib.api_request import api_request
from lib.json_utils import json_utils


def test_change_pass():
    database.reinit()
    requests = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'request.json'))['Requests']

    response = api_request.ApiRequest.do_admin_request(requests[0])
    assert response.status_code== 200
    response = api_request.ApiRequest.do_admin_request(requests[1])
    assert response.status_code== 200
    response = api_request.ApiRequest.do_admin_request(requests[2])
    assert response.status_code== 200