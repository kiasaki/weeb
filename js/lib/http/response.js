class Response {
  constructor(container, response) {
    this.container = container;
    this.response = response;
    this.templateRenderer = container.get("templateRenderer");
  }

  set(headerName, value) {
    this.response.setHeader(headerName, value);
  }

  text(statusCode, value = null) {
    if (!statusCode || typeof statusCode !== "number") {
      throw new Error("Missing statusCode, given: " + statusCode);
    }
    this.statusCode = statusCode;

    const res = this.response;
    if (typeof value === "string") {
      if (!res.getHeader("Content-Type")) {
        res.setHeader("Content-Type", "text/plain; charset=UTF-8");
      }
      res.setHeader("Content-Length", value.length);
      res.writeHead(statusCode);
      res.end(value, "utf8");
    } else {
      res.writeHead(statusCode);
      res.end();
    }
  }

  json(statusCode, value) {
    if (!this.response.getHeader("Content-Type")) {
      this.response.setHeader("Content-Type", "application/json; charset=UTF-8");
    }
    return this.text(statusCode, value);
  }

  noContent(statusCode = 204) {
    return this.text(statusCode);
  }

  async template(statusCode, name, context) {
    const contents = await this.templateRenderer.render(name, context);
    if (!this.response.getHeader("Content-Type")) {
      this.response.setHeader("Content-Type", "text/html; charset=UTF-8");
    }
    return this.text(statusCode, contents);
  }
}

module.exports = Response;
