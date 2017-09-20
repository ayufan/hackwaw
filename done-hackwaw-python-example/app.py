import json
import os
from datetime import datetime, timedelta

import click
import requests
import time
import tinydb
import rfc3339
from flask import Flask, jsonify, Response

TWITTER_URL = os.environ.get("TWITTER_URL", "https://hackwaw-twitter-proxy.herokuapp.com")
SLACK_URL = os.environ.get("SLACK_URL", "https://hackwaw-slack-proxy.herokuapp.com")

LOGO_URL = "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c3/Python-logo-notext.svg/1024px-Python-logo-notext.svg.png"

data_store_path = '/tmp/tweets.json'

print("Creating data store in file {}".format(data_store_path))
store = tinydb.TinyDB(data_store_path).table('tweets')

app = Flask(__name__)
app.config['DEBUG'] = True


@app.route("/health")
def health():
    print("Returning dummy health information")
    status = {
        "app": "OPERATIONAL",
        "twitter": "OPERATIONAL",
        "slack": "OPERATIONAL",
    }

    return jsonify(status)


@app.route('/latest')
def latest():
    print("Loading latest tweets")

    tweets = store.all()
    print("Loaded {} tweets".format(len(tweets)))

    json_output = json.dumps(tweets)
    print("Returning JSON response: {}".format(json_output))

    return Response(json_output, mimetype='application/json')


def load_latest_tweets(since):
    print("Querying twitter proxy for tweets since {}".format(since))

    params = {
        "from": since,
        "to": rfc3339.now().isoformat(),
    }

    response = requests.get(TWITTER_URL + "/tweets", params=params).json()
    print("Twitter proxy responded with {}".format(response))

    return response


def post_on_slack(tweets):

    print("Posting {} tweets to slack".format(len(tweets)))

    for tweet in tweets:
        requests.post(SLACK_URL + "/push", json={
            "icon_url": LOGO_URL,
            "text": tweet['body'],
            "tweetId": tweet['id'],
            "date": tweet["date"],
            "team": "Python team",
        })


# App commands
@click.group()
def commands():
    pass


@commands.command()
def pipe():
    while True:
        print("Loading newest tweets")

        if len(store.all()) > 0:
            since = max(x['date'] for x in store.all())
        else:
            since = (rfc3339.now() - timedelta(seconds=60)).isoformat()

        loaded = load_latest_tweets(since)

        post_on_slack(loaded)
        store.insert_multiple(loaded)

        time.sleep(1)

@commands.command()
def web():
    app.run(port=8080, host='0.0.0.0')


if __name__ == '__main__':
    commands()
