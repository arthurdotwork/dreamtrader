import os

from flask import Flask, request

from stock import get_stock_data

app = Flask(__name__)


@app.route("/stocks/<stock>", methods=["GET"])
def get_stock(stock: str):
    period = request.args.get("period", "1m")
    interval = request.args.get("interval", "1d")

    try:
        stock = get_stock_data(stock, period, interval)
    except Exception as e:
        return {"error": str(e)}, 500

    return {"stock": stock}


def env(name, fallback):
    return os.environ.get(name, fallback)


if __name__ == "__main__":
    from waitress import serve

    try:
        serve(app, host=env("HTTP_HOST", "0.0.0.0"), port=env("HTTP_PORT", 8080))
    except Exception as e:
        print(f"Error: {e}")
        exit(1)
