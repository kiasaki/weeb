const http = require("http");

const Request = require("./request");
const Response = require("./response");

class Server {
  constructor(container, handlerFunc) {
    this.handlerFunc = handlerFunc;
    this.container = container;
    this.events = container.get("events");
    this.logger = container.get("logger");

    this.server = new http.Server(this.handle.bind(this));

    this.server.on("error", err => {
      this.logger.fatal("http server error", { error: err.stack });
      process.exit(1);
    });
  }

  handle(req, res) {
    const request = new Request(this.container, req);
    const response = new Response(this.container, res);
    const startMs = Date.now();
    Promise.resolve().then(() => this.handlerFunc(
      request, response,
    )).catch(err => {
      this.events.emit("http:error", request, response, err);
    }).then(() => {
      request.logger.info("request", {
        method: request.method,
        code: response.statusCode,
        path: request.url.path,
        ms: Date.now() - startMs,
      });
    }).catch(err => {
      request.logger.fatal("fatal error", { error: err.stack });
      response.text(500, "Internal Server Error");
    });
  }

  start(port, host = "0.0.0.0") {
    this.server.listen(port, host, () => {
      const details = this.server.address();

      process.on("SIGTERM", () => {
        this.logger.info("gracefully shutting down");
        this.server.close(process.exit);
      });

      this.logger.info("started", { port: details.port });
    });
  }
}

module.exports = Server;
