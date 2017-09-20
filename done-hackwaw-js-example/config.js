/*jslint node: true */
'use strict';
module.exports = {
    endpoints: {
        tweets: process.env.TWITTER_URL || 'https://hackwaw-twitter-proxy.herokuapp.com',
        slack: process.env.SLACK_URL || 'https://hackwaw-slack-proxy.herokuapp.com'
    },
    tasks: {
        minimalInterval: 1000
    },
    database: {
      filename: '/tmp/tweets.db'
    },
    tweetsPerPage: 15,
    appRoot: __dirname,
    port: process.env.PORT || 8080
};
