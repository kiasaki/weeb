async function globalTemplateStateMiddleware(ctx, next) {
  const config = ctx.container.get("config");
  ctx.state.appName = config.get("appName");
  ctx.state.appBaseUrl = config.get("appBaseUrl");
  await next();
}

export function register(container) {
  const app = container.get("app");
  const template = container.get("services:template");

  app.use(globalTemplateStateMiddleware);

  app.context.render = function(name) {
    this.set("Content-Type", "text/html");
    this.body = template.render(name, this.state);
  };
}
