
/*jslint node: true */
'use strict';

module.exports = function (req, res, next) {
    $log.info('Request:', req.url);
    next();
};
