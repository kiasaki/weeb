const Server = require("http").Server;

function sendText(res, statusCode, value = null) {
  return new Promise(function(resolve) {
    if (typeof value === "string") {
      if (!res.getHeader("Content-Type")) {
        res.setHeader("Content-Type", "text/plain; charset=UTF-8");
      }
      res.setHeader("Content-Length", value.length);
      res.writeHead(statusCode);
      res.end(value, "utf8", resolve);
    } else {
      res.writeHead(statusCode);
      res.end(resolve);
    }
  });
}

function sendJson(res, statusCode, value) {
  if (!res.getHeader("Content-Type")) {
    res.setHeader("Content-Type", "application/json; charset=UTF-8");
  }
  return sendText(res, statusCode, value);
}

function create(fn) {
  const server = new Server((req, res) =>
    new Promise(resolve => resolve(fn(req, res)))
      .then(value => {
        if (value === null) {
          sendText(res, 204, null);
          return;
        }

        if (typeof value === "string") {
          sendText(res, res.statusCode || 200, value);
          return;
        }

        if (value !== undefined) {
          sendJson(res, res.statusCode || 200, value);
        }
      })
      .catch(err => sendError(req, res, err))
  );

  server.on("error", err => {
    console.error("weeb:", err.stack);
    process.exit(1);
  });

  server.run = function(port = process.env.PORT, host = "0.0.0.0") {
    if (!port) port = "3000";
    port = parseInt(port, 10);
    server.listen(port, host, () => {
      const details = server.address();

      process.on("SIGTERM", () => {
        console.log("\nweeb: Gracefully shutting down. Please wait...");
        server.close(process.exit);
      });

      console.log(`weeb: Accepting connections on port ${details.port}`);
    });
  };
}

module.exports = create;
module.exports.sendText = sendText;
module.exports.sendJson = sendJson;
