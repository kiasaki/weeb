const nodemailer = require("nodemailer");
const smtpTransport = require("nodemailer-smtp-transport");

class MailSmtpService {
  constructor(config) {
    this.transport = nodemailer.createTransport(
      smtpTransport({
        host: config.get("mailSmtpHost"),
        port: config.get("mailSmtpPort"),
        auth: {
          user: config.get("mailSmtpUsername"),
          pass: config.get("mailSmtpPassword"),
        },
      }),
    );
  }

  send(payload) {
    return this.transport.sendMail({
      from: `${payload.fromName} <${payload.from}>`,
      to: payload.to,
      subject: payload.subject,
      text: payload.contentPlain,
    });
  }
}

MailSmtpService.dependencyName = "services:mailSmtp";
MailSmtpService.dependencies = ["config"];
module.exports = MailSmtpService;
