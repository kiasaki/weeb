const fs = require("fs");
const url = require("url");
const path = require("path");
const http = require("http");
const { promisify } = require("util");
const {
  hash,
  snakeCase,
  ulid,
  ulidToUuid,
} = require("./utils");

// mailer / http / view / job / data / storage / support

const append = a => b => b.concat([a]);

const contains = item => list => list.indexOf(item) !== -1;

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

// data

class DB {
  constructor(config) {
    this.config = config;
    this.databaseUrl = config.get("databaseUrl");
  }

  query(sql, values) {
    if (!config.get("production")) {
      // TODO use logger
      console.log("executing sql '" + sql + "'");
    }
    return new Promise((resolve, reject) => {
      pg.connect(this.databaseUrl, (err, client, done) => {
        if (err) return reject(err);
        client.query(sql, values, (err, result) => {
          done();
          if (err) return reject(err);
          resolve(result);
        });
      });
    });
  }

  setWheres(query, where = {}) {
    if (Object.keys(where).length === 0) return query;
    query.sql += " where";
    query.values = query.values || [];
    for (let field of Object.keys(where)) {
      const value = where[field];
      if (field.slice(0, 3) === "sql") {
        query.sql += ` ${value[0]}`;
        query.values = query.values.concat(value.slice(1));
      } else {
        query.sql += ` ${field} = ?`;
        query.values.push(value);
      }
    }
    return query;
  }

  async select(table, options = {}) {
    const query = { sql: `SELECT * FROM ${table}`, values: [] };
    this.setWheres(query, options.where);
    if (options.offset) {
      query.sql += " OFFSET ?";
      query.values.push(options.offset);
    }
    if (options.limit) {
      query.sql += " LIMIT ?";
      query.values.push(options.limit);
    }
    if (options.sort) {
      if (!Array.isArray(options.sort) || options.sort.length === 0) {
        throw new Error("db.select: `sort` is not an array");
      }
      let sortClauses = [];
      for (let sort in options.sort) {
        let direction = (sort[0] === "-") ? "DESC" : "ASC";
        if (sort[0] === "-") sort = sort.slice(1);
        sortClauses.push(`? ${direction}`);
        query.values.push(sort);
      }
      query.sql += " ORDER BY " + sortClauses.join(", ");
    }

    const result = await this.query(query.sql, query.values);
    return result.rows;
  }

  async count(table, options) {
    const sql = `SELECT COUNT(*) as count FROM ${table}`;
    const query = { sql, values: [] };
    this.setWheres(query, options.where);
    const result = await this.query(query.sql, query.values);
    return parseInt(result.rows["0"].count, 10);
  }

  selectOne(table, where) {
    return this.select(table, { where, limit: 1 }).then(xs => xs[0]);
  }

  async insert(table, obj) {
    const fields = Object.keys(givenObj);

    const sql = `INSERT INTO ${table}`;
    const query = { sql, values: [] };
    query.sql += ` (${fields.join(", ")})`;
    query.sql += ` VALUES (${fields.map(() => "?").join(", ")})`;
    fields.forEach(field => query.values.push(obj[field]));
    query.sql += " RETURNING *";

    const result = await this.query(query.sql, query.values);
    return result.rows["0"];
  }

  async update(table, where, obj) {
    const fields = Object.keys(obj);

    const sql = `UPDATE ${table} SET `;
    const query = { sql, values: [] };
    const setStatements = [];
    for (let field of fields) {
      if (obj[field] !== undefined && obj[field] !== null) {
        setStatements.push(` ${field} = ?`);
        query.values.push(obj[field]);
      }
    }
    query.sql += setStatements.join(", ");
    this.setWheres(query, where);
    query.sql += " RETURNING *";

    const result = await this.query(query.sql, query.values);
    return result.rows["0"];
  }

  delete(table, where) {
    const query = { sql: `DELETE FROM ${table}`, values: [] };
    this.setWheres(query, where);
    return this.query(query.sql, query.values);
  }

  async save(table, obj, primaryKey = "id") {
    if (obj[primaryKey]) {
      const pk = ulidToUuid(obj[primaryKey]);
      const where = { [primaryKey]: pk };
      delete obj[primaryKey];
      obj = await this.update(table, where, obj);
      obj[primaryKey] = pk;
      return obj;
    }
    obj[primaryKey] = ulidToUuid(ulid());
    return await this.insert(table, obj);
  }
}

class Repository {
  constructor(db) {
    this.db = db;

    this.primaryKey = "id";
    this.table = null;
    this.entityClass = null;
  }

  async find(where, limit, offset) {
    if (!this.entityClass) throw new Error("Entity is missing entityClass");
    const options = {};
    if (where) options.where = where;
    if (limit) options.limit = limit;
    if (offset) options.offset = offset;

    const rows = await this.db.select(this.table, options);

    return rows.map(this.entityClass.fromDatabase.bind(Entity));
  }

  findOne(where) {
    return this.find(where, 1).then(xs => xs[0]);
  }

  findById(pk) {
    return this.findOne({ [this.primaryKey]: pk });
  }

  save(entity) {
    const obj = entity.toDatabase();
    return this.entityClass.fromDatabase(
      this.db.save(table, obj, this.primaryKey)
    );
  }
}

class Entity {
  constructor(params, skipDefaults, givenDefaults = {}) {
    const defaults = skipDefaults ? {} : givenDefaults;
    const obj = Object.assign({}, defaults, params);

    for (let key of Object.keys(obj)) {
      const value = obj[key];
      this[key] = (typeof value === "function") ? value() : value;
    }
  }

  toJson(includePrivate = false) {
    const fieldNames = Object.keys(this.fields || {});
    const privateFields = this.privateFields || [];
    const obj = {};

    for (let key of fieldNames) {
      if (!includePrivate && contains(key, privateFields)) continue;
      obj[key] = this[key];
    }

    return obj;
  }

  toObject() {
    return this.toJson(true);
  }

  toDatabase() {
    const fields = this.fields || {};
    const obj = {};
    for (let key of Object.keys(fields)) {
      const field = fields[key];
      obj[snakeCase(key)] = field.toDatabase(this[key]);
    }
    return obj;
  }

  fromDatabase(obj) {
    const fields = this.fields || {};
    for (let key of Object.keys(fields)) {
      const field = fields[key];
      this[key] = field.fromDatabase(obj[snakeCase(key)]);
    }
    return this;
  }

  static fromDatabase(obj) {
    return new this().fromDatabase(obj);
  }

  static newDate() {
    return new Date();
  }
}

Entity.table = null;

const EntityFields = {
  Text: {
    toDatabase: x => x,
    fromDatabase: x => x,
  },
  Number: {
    toDatabase: x => x,
    fromDatabase: x => x,
  },
  Date: {
    toDatabase: x => x.toISOString(),
    fromDatabase: x => new Date(x),
  },
  Ulid: {
    toDatabase: x => ulidToUuid(x),
    fromDatabase: x => uuidToUlid(x),
  },
  Object: {
    toDatabase: x => JSON.stringify(x),
    fromDatabase: x => JSON.parse(x),
  },
};

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
