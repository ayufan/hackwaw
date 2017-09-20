/*jslint node: true */
'use strict';
GLOBAL.$log = require('../hooks/log.js');
GLOBAL.config = require('../config.js');

const tasks = {
        slack: require(__dirname + '/slack.js'),
        twitter: require(__dirname + '/twitter.js')
    },
    path = require('path'),
    async = require('async');

var interval;

const events = {
    onAppExit: function () {
        //
    }
};

const actions = {
    start: function () {
        $log.service('*** INITIAL CONFIGURATION:', JSON.stringify(config), '***');
        actions.getTwitts();
    },
    getTwitts: function () {
        $log.info('actions -> getTwitts');

        let criteria = {
            from: '2016-03-01T18:54:00+00:00',
            to: '2016-03-31T18:54:00+00:00'
        };

        tasks.twitter.get(criteria, function (error, response, body) {
            if (error) {
                return;
            }

            try {
                body = JSON.parse(body);
            } catch (e) {
                return $log.error('Response body is not a valid JSON');
            }

            if (!(body instanceof Array)) {
                $log.error('Response body is not a JSON array.');
                return;
            }

            var i = 0;

            body.forEach(function (entry) {
                tasks.slack.send(entry, function (error) {
                    if (++i === body.length || body.length === 0) {
                        interval = setTimeout(actions.getTwitts, config.tasks.minimalInterval);
                    }
                });
            });
        });
    }
};

actions.start();

// handle 'exit' events
process.on('exit', events.onAppExit);
process.on('SIGINT', events.onAppExit);
process.on('uncaughtException', events.onAppExit);
// sed '2,4!d' app.js
