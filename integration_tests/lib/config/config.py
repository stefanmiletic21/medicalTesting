import yaml
import os

with open(os.path.join(os.path.dirname(__file__), 'config.yaml')) as f:
    config = yaml.load(f)
