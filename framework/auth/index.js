import Auth from "./auth";

async function loadUserMiddleware(ctx, next) {
  const config = ctx.container.get("config");
  const authProvider = ctx.container.get(config.get("authUserProviderServiceKey"));
  const id = ctx.cookies.get(config.get("authCookieName"), { signed: true });
  if (!id) return await next();

  await authProvider.retriveById(id).then(user => {
    if (!user) return;
    ctx.auth.setUser(user);
    await next();
  });
}

async function requireUserMiddleware(ctx, next) {
  const config = ctx.container.get("config");
  if (!ctx.auth.user) {
    ctx.redirect(config.get("authLoginUrl"));
    return;
  }
  await next();
}

export function register(container) {
  // Set current auth instance in all requests
  container.get("app").use(function(ctx) {
    ctx.auth = new Auth(container, ctx);
  });

  // Load user when user cookie present
  container.get("app").use(loadUserMiddleware);

  // Register "requireUser" middleware in container for use by routes
  container.set("middlewares:auth:requireUser", requireUserMiddleware);
}
