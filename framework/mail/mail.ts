class MailService {
  constructor(
    config,
    templateService,
    consoleMailerService,
    smtpMailerService,
  ) {
    this.config = config;
    this.templateService = templateService;
    this.consoleMailerService = consoleMailerService;
    this.smtpMailerService = smtpMailerService;
  }

  senderImplementation() {
    switch (this.config.get("mailImplementation")) {
      case "smtp":
        return this.smtpMailService;
      case "postmark":
        return this.postmarkMailService;
      default:
        return this.consoleMailService;
    }
  }

  // options = {
  //   from: "",
  //   fromName: "",
  //   to: "",
  //   toName: "",
  //   replyTo: "",
  //   subject: "",
  //   template: "",
  //   templateVars: "",
  // }
  send(options) {
    const { template, templateVars } = options;
    options.contentText = this.templateService.render(
      template + ".txt",
      templateVars,
    );
    options.contentHtml = this.templateService.render(
      template + ".html",
      templateVars,
    );
    options.from = options.from || this.config.get("mailDefaultFrom");
    options.fromName =
      options.fromName || this.config.get("mailDefaultFromName");

    return this.senderImplementation().send(options);
  }
}

MailService.dependencyName = "services:mail";
MailService.dependencies = [
  "config",
  "services:template",
  "services:mailConsole",
  "services:mailPostmark",
  "services:mailSmtp",
];
export default MailService;
