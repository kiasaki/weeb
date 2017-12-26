import * as util from "util";

class Logger {
  context: object;
  parent?: Logger;

  constructor(context = {}, parent?: Logger) {
    this.context = context;
    this.parent = parent;
  }

  addContext(context = {}) {
    this.context = Object.assign(this.context, context);
  }

  withSubContext(context = {}) {
    return new Logger(context, this);
  }

  debug(message: string, extraArgs: object) {
    this.log("debug", message, extraArgs);
  }

  info(message: string, extraArgs: object) {
    this.log("info", message, extraArgs);
  }

  warning(message, extraArgs) {
    this.log("warning", message, extraArgs);
  }

  error(message, extraArgs) {
    this.log("error", message, extraArgs);
  }

  fatal(message, extraArgs) {
    this.log("fatal", message, extraArgs);
  }

  log(level, message, extraArgs = {}) {
    const time = new Date().toISOString();

    // Loop over parent to contruct final "context"
    let context = {};
    let subject = this;
    while (subject) {
      context = Object.assign(subject.context, context);
      subject = subject.parent;
    }

    const logMessage = Object.assign(context, extraArgs, {
      time,
      level,
      message,
    });

    if (process.stdout.isTTY) {
      const formattedMessage = util
        .inspect(logMessage, { colors: true })
        .replace(/\n/g, "")
        .replace(/ {2}/g, " ")
        .replace(/{ /g, "{")
        .replace(/ }/g, "}");
      console.log(formattedMessage);
    } else {
      console.log(JSON.stringify(logMessage));
    }
  }
}

export default Logger;
