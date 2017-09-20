/*jslint node: true */
'use strict';
var request = require('request');
var path = require('path');
var _ = require('lodash');
var colors = require('colors');
var Datastore = require('nedb');
var db = new Datastore({ filename: path.resolve([__dirname, config.database.filename].join('/../')) });

db.loadDatabase(function (err) {
    if (err) {
        throw new Error(err);
    }
});

module.exports = {
    send: function (inputData, callback) {

        var data = {
            date: inputData.date,
            text: inputData.body,
            tweetId: inputData.id
        };

        // default values to be posted
        var defaults = {
            date: '2016-03-01T18:00:00Z',
            icon_url: 'http://www.veryicon.com/icon/ico/System/Arcade%20Daze/Mario.ico',
            team: 'NodeJS_BOT',
            text: 'Hello, World!',
            tweetId: 0
        };

        data = _.defaults(data, defaults);

        db.insert(inputData, function (err, newDocs) {
            if (err) {
                $log.error('DB Error');
            }
        });

        request.post(
            {
                url: config.endpoints.slack + "/push",
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            },
            function(err, httpResponse, body) {
                if (err) {
                    return $log.error('[slack]', err);
                }

                $log.debug('[slack] Sent:', JSON.stringify(data), 'Status code:', httpResponse.statusCode);
                callback.apply(null, arguments);
            }
        );
    }
};
