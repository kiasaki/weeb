const path = require("path");

const Config = require("./config");
const Container = require("./container");
const Events = require("./events");
const Logger = require("./logger");
const Router = require("./http/router");
const Server = require("./http/server");
const TemplateRenderer = require("./view/template_renderer");
const { append } = require("./util");
const { MIME_TYPES } = require("./constants");
const { HttpError, ServerErrorHttpError } = require("./errors");

class Project {
  constructor() {
    this.container = new Container();
    this.container.set("application", this);
    this.container.set("container", this.container);

    this.events = new Events();
    this.container.set("events", this.events);

    this.config = new Config();
    this.container.set("config", this.config);

    this.logger = new Logger();
    this.container.set("logger", this.logger);

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

    // Setup default error handler
    this.events.on("http:error", (req, res, err) => {
      const originalError = err;
      if (!(err instanceof HttpError)) {
        err = new ServerErrorHttpError();
      }
      if (err.statusCode >= 500) {
        req.logger.error("server error", { error: originalError.stack });
      }
      res.text(err.statusCode, err.message);
    })

    // Setup routes
    this.router.get("/static/(.*)", this.handleStatic.bind(this));
  }

  async handleStatic(req, res) {
    const cacheFiles = this.config.get("production");
    const fs = new FileSystem(cacheFiles);
    const folders = this.config.get("staticFolders");
    const filePathName = req.url.pathname.slice("/static".length).replace('..', '');
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

module.exports = Project;
