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
        px = str(trade['LastPx'])
        if '.' in px:
            a,b = px.split('.')
            if not ((len(a) == 3 and len(b) == 6) or
                    (len(a) == 4 and len(b) == 5) or
                    (len(a) == 5 and len(b) == 4) or
                    (len(a) == 6 and len(b) == 3)):
                print('{}: {}v{}'.format(px, len(a), len(b)))
                valid = False
                break

    response = {'statusCode': 200, 'body': {'valid':valid, 'portfolio': portfolio}}
    endTime = 1000*time.time()
    return json.dumps(timestamp(response, event,startTime, endTime, 0)), 200
