/*jslint node: true */
'use strict';
var request = require('request');
var _ = require('lodash');
var colors = require('colors');

module.exports = {
    $TAG: '[twitter]',
    get: function (criteria, callback) {

        const _self = this;

        let defaults = {
            from: '2016-03-01T18:00:00.000',
            to: '2016-03-31T18:00:00.000'
        };

        criteria = _.defaults(criteria, defaults);

        $log.info(_self.$TAG, 'Getting latest twitts with criteria:', JSON.stringify(criteria));

        request.get(
            {
                url: config.endpoints.tweets + "/tweets",
                headers: {
                    'Content-Type': 'application/json'
                },
                qs: criteria
            },
            function(error, response, body) {
                if (error) {
                    return $log.error(_self.$TAG, error);
                }

                $log.debug(_self.$TAG, 'Respone status code is:', response.statusCode);
                callback.apply(null, arguments);
            }
        );

    }
};
