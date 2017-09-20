"use strict";

GLOBAL.$log = require([__dirname, '/hooks/log.js'].join(''));
GLOBAL.config = require([__dirname, '/config.js'].join(''));

var handlers = {
    request: require([__dirname, '/hooks/request-handler.js'].join('')),
    notFound: require('./hooks/404notFound.js')
};

var swagger = require('swagger-express-mw');
var path = require('path');
var app = require('express')();
var colors = require('colors');
var Datastore = require('nedb');
GLOBAL.db = new Datastore({ filename: config.database.filename });

db.loadDatabase(function (err) {
    if (err) {
        throw new Error(err);
    }
});

module.exports = app;

app.use(handlers.request);

// create swagger app
swagger.create(config, function (err, swaggerExpress) {
    if (err) {
        throw err;
    }

    // register swagger routes
    swaggerExpress.register(app);

    // handle 404 not found
    app.get('*', handlers.notFound);

    // start http listener on given port
    app.listen(config.port);

    $log.service(['Go to: http://localhost:', config.port, '/health'].join(''))
});
