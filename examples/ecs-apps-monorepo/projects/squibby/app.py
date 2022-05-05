# flask_web/app.py
import os
import ddtrace
from flask import Flask
from ddtrace import tracer

app = Flask(__name__)
appname = os.getenv('APP_NAME', 'Masyanya')
@app.route('/')
def hello_world():
    return appname


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=3000)
