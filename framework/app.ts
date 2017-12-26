import * as http from "http";
import * as Koa from "koa";
import * as KoaRouter from "koa-router";

import Container from "./container";

async function requestLoggerMiddleware(ctx: Koa.Context, next: Function) {
  const logger = ctx.container.get("logger");
  const start = Date.now();
  await next();
  const ms = Date.now() - start;
  logger.info("request", { method: ctx.method, url: ctx.url, ms });
}

async function forceSslMiddleware(ctx: Koa.Context, next: Function) {
  const config = ctx.container.get("config");

  if (!ctx.secure && config.get("appForceSsl")) {
    ctx.redirect("https://" + config.get("appBaseUrl") + ctx.path);
    return;
  }

  await next();
}

export function preRegister(container: Container) {
  const logger = container.get("logger");
  const config = container.get("config");

  const app = container.set("app", new Koa());
  const router = container.set("router", new KoaRouter());
  app.context.container = container;
  app.context.router = router;

  // Pre-Request Middlewares
  app.proxy = true;
  app.use(forceSslMiddleware);
  app.use(requestLoggerMiddleware);

  app.start = function() {
    container.get("events").trigger("app:start");
    process.on("SIGINT", function() {
      logger.info("SIGINT caught");
      app.stop();
      process.exit();
    });
    app.server = http.createServer(app.callback());
    app.server.listen(config.get("port"), () => {
      logger.info("app started", { port: config.get("port") });
      container.get("events").trigger("app:started");
    });
  };

  app.stop = function() {
    logger.info("app stopping");
    app.server.close(() => {
      logger.info("app stopped");
      container.get("events").trigger("app:stopped");
    });
  };
}

export function postRegister(container: Container) {
  const app = container.get("app");
  const config = container.get("config");
  const router = container.get("router");

  app.keys = config.get("authCookieKeys").split(",");

  // Post-Request Middlewares
  app.use(router.routes());
  app.use(router.allowedMethods());
  app.use(async function(ctx: Koa.Context) {
    if (ctx.accepts("html", "json") === "json") {
      ctx.body = { error: "not found" };
    } else {
      await ctx.render(config.get("appNotFoundTemplate"));
    }
  });
}
