/*jslint node: true */
'use strict';
var util = require('util');

module.exports = {
  latest: getLatest
};

function getLatest(req, res) {
    let page = req.swagger.params.page.value || 0;

    db.loadDatabase();
    db.find({})
      .sort({ id: -1 })
      .skip(page * config.tweetsPerPage)
      .limit(15)
      .exec(function (err, tweets) {
          if (err) {
              throw new Error(err);
          }

          res.json(tweets);
      });
}
