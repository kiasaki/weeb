const http = require("http");

// mailer / http / view / job / data / storage / support

class Container {
  constructor() {
    this.values = {};
  }

  get(key) {
    return this.values[key];
  }

  set(key, value) {
    return (this.values[key] = value);
  }
}

class Request {
  constructor(container, request) {
    this.container = container;
    this.request = request;

    this.url = request.url;
    this.method = request.method;
    this.params = {}; // Filled in by router
  }
}

class Response {
  constructor(container, response) {
    this.container = container;
    this.response = response;
  }

  text(statusCode, value = null) {
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
    this.text(statusCode, value);
  }

  noContent(statusCode = 204) {
    this.text(statusCode);
  }
}

class Router {
  constructor() {
    this.handle = this.handle.bind(this);
    this.routes = {
      HEAD: [],
      OPTIONS: [],
      GET: [],
      POST: [],
      PUT: [],
      DELETE: [],
    };
    this.notFoundHandler = (req, res) => {
      res.text(404, "Not Found");
    }
  }

  add(method, path, handler) {
    const paramRe = new RegExp("<([a-zA-Z]+)>", "g");
    const paramNames = [];
    let param = paramRe.exec(path);
    while (param) {
      paramNames.push(param[1]);
      param = paramRe.exec(path);
    }
    this.routes[method].push({
      path: path,
      paramNames: paramNames,
      handler: handler,
    });
  }

  setNotFound(handler) {
    this.notFoundHandler = handler;
  }

  handle(request, response) {
    const paramRe = new RegExp("<([a-zA-Z]+)>", "g");
    const routes = this.routes[request.method];

    for (let i = 0; i < routes.length; i++) {
      const route = routes[i];
      const result = new RegExp(route.path.replace(paramRe, "([^/]+)")).exec(request.url);
      if (result) {
        const paramValues = result.slice(1);
        for (let i = 0; i < route.paramNames.length; i++) {
          request.params[route.paramNames[i]] = paramValues[i];
        }
        router.handler(request, response);
        return;
      }
    }

    // Didn't match any route, render 404 / Not Found
    this.notFoundHandler(request, response);
  }
}

class Server {
  constructor(container, handlerFunc) {
    this.handlerFunc = handlerFunc;
    this.container = container;
    this.server = new http.Server(this.handle.bind(this));
    this.server.on("error", err => {
      console.error("weeb:", err.stack);
      process.exit(1);
    });
  }

  handle(req, res) {
    this.handlerFunc(
      new Request(this.container, req),
      new Response(this.container, res),
    );
  }

  start(port, host = "0.0.0.0") {
    this.server.listen(port, host, () => {
      const details = this.server.address();

      process.on("SIGTERM", () => {
        console.log("\nweeb: Gracefully shutting down. Please wait...");
        this.server.close(process.exit);
      });

      console.log(`weeb: Accepting connections on port ${details.port}`);
    });
  }
}

class Application {
  constructor() {
    this.container = new Container();
    this.container.set("application", this);

    this.router = new Router();
    this.container.set("router", this.router);

    this.server = new Server(this.container, this.router.handle);
    this.container.set("server", this.server);
  }

  start(port, host = "0.0.0.0") {
    port = parseInt(port || process.env.PORT || "3000", 10);
    this.server.start(port, host);
  }
}

module.exports.Application = Application;
module.exports.Server = Server;
module.exports.Request = Request;
module.exports.Response = Response;
