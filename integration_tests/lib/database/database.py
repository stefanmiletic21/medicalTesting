import psycopg2
from lib.config.config import config
from pathlib2 import Path
import os

params = {
    'database': config['database']['database'],
    'user': config['database']['user'],
    'host': config['database']['host'],
    'password': config['database']['password'],
    'port': config['database']['port']
}

def query(string, *args):
    conn = psycopg2.connect(**params)
    cur = conn.cursor()

    cur.execute(string, args)
    result = cur.fetchall()

    return result

def execute(string, *args):
    conn = psycopg2.connect(**params)
    cur = conn.cursor()

    cur.execute(string, args)
    conn.commit()

    rowcount = cur.rowcount

    cur.close()
    return rowcount

def clean():
    conn = psycopg2.connect(**params)
    cur = conn.cursor()

    contents = Path(os.path.join(os.path.dirname(__file__), 'clean.sql')).read_text()
    cur.execute(contents)
    conn.commit()

    rowcount = cur.rowcount

    cur.close()
    return rowcount

def reinit():
    conn = psycopg2.connect(**params)
    cur = conn.cursor()

    contents = Path(os.path.join(os.path.dirname(__file__), 'clean.sql')).read_text()
    cur.execute(contents)
    conn.commit()

    contents = Path(os.path.join(os.path.dirname(__file__), 'dummy_data.sql')).read_text()
    cur.execute(contents)
    conn.commit()

    cur.close()
    return 