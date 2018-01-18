import * as Koa from "koa";

import Container from "./container";

declare interface ErrorWithStatus extends Error {
  status?: number;
}

export function register(container: Container) {
  const app = container.get("app");
  const config = container.get("config");
  const logger = container.get("logger");

  app.on("error", async (err: ErrorWithStatus, ctx: Koa.Context) => {
    logger.error("server error", { error: err.message, stack: err.stack });
    ctx.status = err.status || 500;

    if (ctx.accepts("json")) {
      ctx.body = { error: err.message };
    } else {
      ctx.state.status = err.status || 500;
      ctx.state.message = err.message;
      const env = config.get("nodeEnv");
      if (env !== "production") {
        ctx.state.error = err;
      }
      await ctx.render(config.get("appErrorTemplate"));
    }
  });
}
