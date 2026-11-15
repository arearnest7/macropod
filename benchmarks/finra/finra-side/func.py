import time
import json
from util import *

def FunctionHandler(context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    event = context["JSON"]

    startTime = 1000*time.time()

    portfolio = event['body']['portfolio']
    portfolios = json.loads(open('portfolios.json', 'r').read())
    data = portfolios[portfolio]

    valid = True

    for trade in data:
        side = trade['Side']
        # Tag ID: 552, Tag Name: Side, Valid values: 1,2,8
        if not (side == 1 or side == 2 or side == 8):
            valid = False
            break

    response = {'statusCode': 200, 'body': {'valid':valid, 'portfolio': portfolio}}
    endTime = 1000*time.time()
    return json.dumps(timestamp(response, event, startTime, endTime, 0)), 200
