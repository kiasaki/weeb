const { NotFoundHttpError } = require("../errors");

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

    this.head = this.add.bind(this, "HEAD");
    this.options = this.add.bind(this, "OPTIONS");
    this.get = this.add.bind(this, "GET");
    this.post = this.add.bind(this, "POST");
    this.put = this.add.bind(this, "PUT");
    this.delete = this.add.bind(this, "DELETE");
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

  handle(request, response) {
    const paramRe = new RegExp("<([a-zA-Z]+)>", "g");
    const routes = this.routes[request.method];

    for (let i = 0; i < routes.length; i++) {
      const route = routes[i];
      const result = new RegExp(route.path.replace(paramRe, "([^/]+)")).exec(request.url.pathname);
      if (result) {
        const paramValues = result.slice(1);
        for (let i = 0; i < route.paramNames.length; i++) {
          request.params[route.paramNames[i]] = paramValues[i];
        }
        return route.handler(request, response);
      }
    }

    // Didn't match any route
    throw new NotFoundHttpError();
  }
}

module.exports = Router;
