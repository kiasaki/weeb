const PREFIX = "/static/";
const PREFIX_LENGTH = PREFIX.length;

async function staticMiddleware(ctx, next) {
  if (ctx.method !== "HEAD" && ctx.method !== "GET") return await next();
  if (ctx.path.slice(0, PREFIX_LENGTH) !== PREFIX) return await next();

  const config = ctx.container.get("config");
  for (folder of config.get("appStaticFolders")) {
    if (await send(ctx, ctx.path.slice(PREFIX_LENGTH), { root: folder })) {
      return;
    }
  }

  await next();
}

export function preRegister(container) {
  const app = container.get("app");
  const config = container.get("config");

  config.set("appStaticFolders", []);

  app.addStaticFolder = function(folderPath) {
    const folders = config.get("appStaticFolders");
    config.set("appStaticFolders", folders.concat([folderPath]));
  };

  app.use(staticMiddleware);
}
