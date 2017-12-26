'use strict';

const User = require('../entities/user');
const crypto = require('../helpers/crypto');

const TWENTY_FOUR_HOURS_IN_MINUTES = 24 * 60;

class PasswordResetService {
  constructor(i18n, userRepository, userService, jwtService, mailerService) {
    this.i18n = i18n;
    this.userRepository = userRepository;
    this.userService = userService;
    this.jwtService = jwtService;
    this.mailerService = mailerService;
    this.generatePassword = crypto.generatePassword;

    this.forgot = this.forgot.bind(this);
  }

  forgot(email) {
    return this.userRepository.findByEmail(email)
      .then(user => {
        const token = this.generateToken(user.id);
        return this.sendResetInstructions(user.fname, user.email, token);
      });
  }

  generateToken(id) {
    return this.jwtService.sign({
      subject: 'password-reset', id
    }, TWENTY_FOUR_HOURS_IN_MINUTES);
  }

  sendResetInstructions(name, email, token) {
    return this.mailerService.send({
      to: email,
      subject: this.i18n.t('app.services.password-reset.sendResetInstructions.subject'),
      template: 'emails-forgot',
      templateVars: {name, token}
    });
  }

  reset(token, newPassword) {
    return this.jwtService.verify(token)
      .then(tokenContents => {
        if (tokenContents.subject !== 'password-reset') {
          throw new Error('Invalid password reset token');
        }
        return this.userRepository.findById(tokenContents.id);
      })
      .then(user => {
        return Promise.all([
          user,
          this.generatePassword(newPassword)
        ]);
      })
      .then(args => {
        const user = args[0];
        const hashedPassword = args[1];
        user.password = hashedPassword;
        return this.userRepository.save(new User(user));
      });
  }
}

PasswordResetService.dependencyName = 'services:password-reset';
PasswordResetService.dependencies = [
  'i18n', 'repositories:user', 'services:user',
  'services:jwt', 'services:mailer'
];
module.exports = PasswordResetService;
