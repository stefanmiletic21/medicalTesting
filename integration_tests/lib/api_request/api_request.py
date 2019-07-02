import json

from lib.header import header
from lib.config.config import config
from enum import Enum
import requests


class ApiRequest(object):
    def __init__(self, header, api, data, method):
        self.header = header
        self.request_method = RequestMethod[method]
        self.data = data
        self.api = config['api_url'] + api
    def do_request(self):
        if self.request_method == RequestMethod.POST:
            return requests.post(self.api, data=self.data, headers=self.header)
        elif self.request_method == RequestMethod.GET:
            return requests.get(self.api, headers=self.header)
        elif self.request_method == RequestMethod.PUT:
            return requests.put(self.api, data=self.data, headers=self.header)
        elif self.request_method == RequestMethod.DELETE:
            return requests.delete(self.api, headers=self.header)
        elif self.request_method == RequestMethod.PATCH:
            return requests.patch(self.api, data=self.data, headers=self.header)
        
    @staticmethod
    def login_test_request(request, auth=""):
        head = header.get_admin_header()
        if auth != "":
            head["Authorization"] = auth
        response = ApiRequest(head, "testingHash", "{}", "POST").do_request()
        responseContent = json.loads(response.content)
        head["Authorization"] = responseContent["Auth"]
        
        return ApiRequest(head, request['API'], json.dumps(request['Data']), request['Method']).do_request()
        
    @staticmethod
    def do_admin_request(request):
        head = header.get_admin_header()
        response = ApiRequest(head, "/testingHash", "{}", "POST").do_request()
        responseContent = json.loads(response.content)
        head["Authorization"] = responseContent["Auth"]
        
        return ApiRequest(head, request['API'], json.dumps(request['Data']), request['Method']).do_request()

    @staticmethod
    def do_nurse_request(request):
        head = header.get_nurse_header()
        response = ApiRequest(head, "/testingHash", "{}", "POST").do_request()
        responseContent = json.loads(response.content)
        head["Authorization"] = responseContent["Auth"]
        
        return ApiRequest(head, request['API'], json.dumps(request['Data']), request['Method']).do_request()
    @staticmethod
    def do_doctor_request(request):
        head = header.get_doctor_header()
        response = ApiRequest(head, "/testingHash", "{}", "POST").do_request()
        responseContent = json.loads(response.content)
        head["Authorization"] = responseContent["Auth"]
        
        return ApiRequest(head, request['API'], json.dumps(request['Data']), request['Method']).do_request()

class RequestMethod(Enum):
    POST = 1
    GET = 2
    PUT = 3
    DELETE = 4
    PATCH = 5