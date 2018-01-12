const fs = require("fs");
const url = require("url");
const path = require("path");
const http = require("http");
const { promisify } = require("util");
const { hash } = require("./utils");

// mailer / http / view / job / data / storage / support

const append = a => b => b.concat([a]);

const MIME_TYPES = {
  "js": "text/javascript",
  "ico": "image/x-icon",
  "css": "text/css",
  "png": "image/png",
  "jpg": "image/jpeg",
  "wav": "audio/wav",
  "mp3": "audio/mpeg",
  "svg": "image/svg+xml",
  "pdf": "application/pdf",
  "doc": "application/msword",
  "eot": "appliaction/vnd.ms-fontobject",
  "ttf": "aplication/font-sfnt",
  "html": "text/html",
  "json": "application/json",
};

// support

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

class Config {
  constructor() {
    this.values = {};
  }

  get(key, default_) {
    return this.values[key] || default_;
  }

  set(key, value) {
    return (this.values[key] = value);
  }

  update(key, updateFn) {
    this.values[key] = updateFn(this.values[key]);
  }
}

class FileSystem {
  constructor(cached = false) {
    this.cached = cached;
    this.fileCache = {};
  }

  async findInDirectories(name, folders) {
    for (let folder of folders) {
      try {
        const fullPath = path.join(folder, name);
        await promisify(fs.stat)(fullPath);
        return fullPath;
      } catch (err) {
        if (err.code === "ENOENT") {
          continue;
        }
        throw err;
      }
    }
    return null;
  }

  async isDirectory(name) {
    const stat = await promisify(fs.stat)(name);
    return stat.isDirectory();
  }

  async read(name) {
    if (name in this.fileCache) return this.fileCache[name];

    const contents = await promisify(fs.readFile)(name, { encoding: "utf8" });
    this.fileCache[name] = contents;
    return contents;
  }
}

// view

class TemplateRenderer {
  constructor(config) {
    this.config = config;
    this.templateCache = {};
    const cacheFiles = config.get("production");
    this.fs = new FileSystem(cacheFiles);
  }

  compile(templateSource) {
    const code = "var p=[],print=function(){p.push.apply(p,arguments);};" +
      "with(context||{}){p.push('" +
      templateSource
        .replace(/[\r\t\n]/g, " ")
        .split("<%").join("\t")
        .replace(/((^|%>)[^\t]*)'/g, "$1\r")
        .replace(/\t=(.*?)%>/g, "',$1,'")
        .split("\t").join("');")
        .split("%>").join("p.push('")
        .split("\r").join("\\'")
      + "');}return p.join('');";
    return new Function("context", code);
  }

  async render(name, context) {
    const directories = this.config.get("templateFolders");
    const path = await this.fs.findInDirectories(name + ".html", directories);
    if (!path) {
      throw new Error("Template named '" + name + "' not found");
    }
    const templateSource = await this.fs.read(path);
    const templateSourceHash = hash(templateSource);
    const template = this.templateCache[templateSourceHash] || this.compile(templateSource);
    if (this.config.get("production")) {
      this.templateCache[templateSourceHash] = template;
    }
    return template(context);
  }
}

// http

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
    this.templateRenderer = container.get("templateRenderer");
  }

  set(headerName, value) {
    this.response.setHeader(headerName, value);
  }

  text(statusCode, value = null) {
    if (!statusCode || typeof statusCode !== "number") {
      throw new Error("Missing statusCode, given: " + statusCode);
    }

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
        return route.handler(request, response);
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
    const request = new Request(this.container, req);
    const response = new Response(this.container, res);
    Promise.resolve().then(() => this.handlerFunc(
      request, response,
    )).catch(err => {
      console.log("weeb:", err.stack);
      // TODO custom error handler
      response.text(500, "Internal Server Error");
    });
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
    this.container.set("container", this.container);

    this.config = new Config();
    this.container.set("config", this.config);

    this.router = new Router();
    this.container.set("router", this.router);

    this.server = new Server(this.container, this.router.handle);
    this.container.set("server", this.server);

    this.container.set("templateRenderer", new TemplateRenderer(this.config));

    // Setup config
    this.config.set("production", process.env.NODE_ENV === "production");
    this.config.set("apps", []);
    this.config.set("templateFolders", []);
    this.config.set("staticFolders", []);

    // Setup routes
    this.router.get("/static/(.*)", this.handleStatic.bind(this));
  }

  async handleStatic(req, res) {
    const cacheFiles = this.config.get("production");
    const fs = new FileSystem(cacheFiles);
    const folders = this.config.get("staticFolders");
    const filePathName = url.parse(req.url).pathname.slice("/static".length);
    let filePath = await fs.findInDirectories(filePathName, folders);
    if (!filePath) {
      res.text(404, "File not found");
      return;
    }

    if (await fs.isDirectory(filePath)) {
      filePath = path.join(filePath, "index.html");
    }

    const contents = await fs.read(filePath);
    const ext = path.parse(filePath).ext.slice(1);
    res.set("Content-Type", MIME_TYPES[ext] || "text/plain");
    res.text(200, contents);
  }

  addApplication(name, rootFolder) {
    this.config.update("apps", append({ name, rootFolder }));
    this.config.update("templateFolders", append(path.join(rootFolder, "templates")));
    this.config.update("staticFolders", append(path.join(rootFolder, "static")));

    try {
      require(path.join(rootFolder, "routes.js"))(this.router, this.container);
    } catch(e) { /* ignore */ }

    // TODO load controllers / services / routes and more
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
