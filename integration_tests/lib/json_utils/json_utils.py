import json

def parse_json(file):
    i = json.loads(read_json(file))
    return i


def read_json(file):
    i = open(file).read()
    return i


def dump_json(data, file):
    with open(file, "w") as jsonFile:
        i = json.dump(data, jsonFile, indent=2)
        return i