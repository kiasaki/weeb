import postmark from "postmark";
import { promisify } from "util";

class MailPostmarkService {
  constructor(config) {
    const apiToken = config.get("mailPostmarkApiToken");
    this.client = new postmark.Client(apiToken);
    this.sendPostmarkEmail = promisify(this.client.sendEmail.bind(this.client));
  }

  send(spec) {
    const email = {
      From: spec.from,
      To: spec.to,
      Subject: spec.subject,
      TextBody: spec.contentPlain,
      HtmlBody: spec.contentHtml,
    };

    if (spec.replyTo) {
      email.ReplyTo = spec.replyTo;
    }

    return this.sendPostmarkEmail(email);
  }
}

MailPostmarkService.dependencyName = "services:mailPostmark";
MailPostmarkService.dependencies = ["config"];
export default MailPostmarkService;
