# flask_web/app.py
import os
import ddtrace
from flask import Flask
from ddtrace import tracer

app = Flask(__name__)
api_key = os.getenv('EXAMPLE_API_KEY')
secret_key = os.getenv('EXAMPLE_SECRET')
text = '👺 This is api key: ' + api_key + '! This is secret: '+ secret_key + ' 👺'
@app.route('/')
def hello_world():
    return (text)

@app.route('/health')
def health():
    return (text)


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=3000)
