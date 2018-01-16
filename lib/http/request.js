const url = require("url");

class Request {
  constructor(container, request) {
    this.container = container;
    this.request = request;

    this.url = url.parse(request.url);
    this.method = request.method;
    this.params = {}; // Filled in by router

    this.logger = container.get("logger").withSubContext({
      path: this.url.path,
    });
  }
}

module.exports = Request;
