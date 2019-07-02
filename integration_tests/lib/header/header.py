import os
from lib import database
from lib.json_utils import json_utils

def get_admin_header():
    header = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'admin.json'))
    return header

def get_doctor_header():
    header = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'doctor.json'))
    return header

def get_nurse_header():
    header = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'nurse.json'))
    return header

def get_researcher_header():
    header = json_utils.parse_json(os.path.join(os.path.dirname(__file__), 'researcher.json'))
    return header