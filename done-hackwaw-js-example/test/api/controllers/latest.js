'use strict';

var should = require('should');
var request = require('supertest');
var server = require('../../../app');

describe('controllers', function () {

    describe('latest', function () {

        describe('GET /latest', function() {

            it('Should not response with error', function(done) {
                request(server)
                    .get('/latest')
                    .set('Accept', 'application/json')
                    .expect('Content-Type', /json/)
                    .expect(200)
                    .end(function(err, res) {
                        should.not.exist(err);
                        done();
                    });
            });

            it('Response body should be an instance of Array', function(done) {
              request(server)
                  .get('/latest')
                  .set('Accept', 'application/json')
                  .expect('Content-Type', /json/)
                  .expect(200)
                  .end(function(err, res) {
                      res.body.should.be.an.instanceOf(Array);
                      done();
                  });
            });

        });

    });

});
