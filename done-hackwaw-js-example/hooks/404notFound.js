/*jslint node: true */
'use strict';
var path = require('path');

module.exports = function (req, res) {
    $log.error('404 Not found', [' -> original request: ', req.url].join('').grey.reset);
    res.sendFile(path.resolve(__dirname+'/../templates/404notFound.html'));
};
