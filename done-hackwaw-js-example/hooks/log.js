/*jslint node: true */
'use strict';

var beeper = require('beeper');

module.exports = {
    info: function (str) {
        console.info([new Date().toString().grey, ' | '.grey, 'INFO: '.blue, Array.prototype.slice.call(arguments).join(' ').blue].join(''));
    },
    debug: function (str) {
        console.info([new Date().toString().grey, ' | '.grey, 'DEBUG: '.green, Array.prototype.slice.call(arguments).join(' ').green].join(''));
    },
    warn: function (str) {
        console.warn([new Date().toString().grey, ' | '.grey, 'WARNING: '.yellow, Array.prototype.slice.call(arguments).join(' ').yellow].join(''));
    },
    error: function (str) {
        console.error([new Date().toString().grey, ' | '.grey, 'ERROR: '.white.bgRed, Array.prototype.slice.call(arguments).join(' ').white.bgRed].join(''));
        beeper();
    },
    service: function (str) {
        console.error([new Date().toString().grey, ' | '.grey, Array.prototype.slice.call(arguments).join(' ').cyan].join(''));
    }
};
