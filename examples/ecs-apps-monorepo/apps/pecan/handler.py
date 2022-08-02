import json
import requests

def endpoint(event, context):
    # get `usd_amount` parameter from event in serverless
    amount =  json.loads(event["queryStringParameters"]["usd_amount"], parse_float=float)
    # get rates of currencies to usd from website
    get_rates = requests.get('http://www.floatrates.com/daily/usd.json')
    # deserialize JSON data to Python variable
    rdata = json.loads(get_rates.text)
    # form the body of the message with data we have from previous iterations
    body = {
        "USD": amount,
        "EUR": rdata['eur']['rate']*amount,
        "RUB": rdata['rub']['rate']*amount,
        "BYN": rdata['byn']['rate']*amount,
    }

    response = {
        "statusCode": 200,
        "body": json.dumps(body)
    }

    return response
