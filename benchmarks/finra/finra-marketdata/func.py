import pandas as pd
import time
import json

def FunctionHandler(context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    event = context["JSON"]

    startTime = time.time()
    externalServicesTime = []
    portfolioType = event['body']['portfolioType']

    tickersForPortfolioTypes = {'S&P': ['GOOG', 'AMZN', 'MSFT']}
    tickers = tickersForPortfolioTypes[portfolioType]

    prices = {}
    whole_set = pd.read_csv("yfinance.csv")
    for ticker in tickers:
        # tickerObj = base.Ticker(ticker)
        #Get last closing price
        tickTime = time.time()
        # data = tickerObj.history(period="1")
        externalServicesTime.append(time.time() - tickTime)
        # price = data['Close'].unique()[0]
        price = whole_set['Close'].unique()[0]
        prices[ticker] = price

    # prices = {'GOOG': 1732.38, 'AMZN': 3185.27, 'MSFT': 221.02}

    endTime = time.time()

    response = {'time': {'start': startTime, 'end': endTime, 'externalServicesTime': externalServicesTime}, 'body': {'marketData':prices}}

    return json.dumps(response), 200
