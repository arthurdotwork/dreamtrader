import pytz
import yfinance


def get_stock_data(stock: str, history: str = "1m", interval: str = "1d") -> object:
    """
    get_stock_data fetches stock data from Yahoo Finance
    based on the stock ticker, history, and interval.
    """
    data = yfinance.Ticker(stock)

    history = data.history(period=history, interval=interval)
    price_data = history.to_dict()['Close']
    if not price_data:
        raise Exception("stock not found")

    stock_data = []
    for key, value in price_data.items():
        timestamp_utc = key.astimezone(pytz.UTC)
        timestamp_str = timestamp_utc.strftime('%Y-%m-%d %H:%M:%S%z')
        stock_data.append({
            "price": value,
            "timestamp": timestamp_str,
        })

    return {
        "ticker": stock,
        "prices": stock_data,
    }
