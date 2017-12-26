'use strict';

const jwt = require('jsonwebtoken');
const promisify = require('../helpers/promisify');

class JwtService {
  constructor(config) {
    this.config = config;
  }

  sign(value, minutesBeforeExpiration) {
    return jwt.sign(value, this.config.get('jwt_secret'), {
      algorithm: 'HS256',
      expiresIn: minutesBeforeExpiration * 60
    });
  }

  verify(token) {
    return promisify(jwt.verify.bind(jwt))(
      token, this.config.get('jwt_secret'), {
        algorithms: ['HS256']
      });
  }
}
JwtService.dependencyName = 'services:jwt';
JwtService.dependencies = ['config'];
module.exports = JwtService;
