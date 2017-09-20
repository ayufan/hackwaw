/*jslint node: true */
'use strict';
var util = require('util');

module.exports = {
  health: getHealth
};


function getHealth(req, res) {
    res.json({
        app: 'OPERATIONAL',
        database: 'OPERATIONAL'
    });
}
