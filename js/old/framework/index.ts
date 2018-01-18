import Config from "./config";
import Container from "./container";
import Logger from "./logger";
import observable from "./observable";

// Takes an empty container and bootstraps core services and app providers into it
export function bootstrap(container: Container, userProviders = []) {
  container.reset();
  container.set("container", container);

  // Setup events
  container.set("events", observable({}));

  // Setup logger
  const logger = new Logger();
  container.setLogger(logger);
  container.set("logger", logger);
  logger.info("start", { component: "bootstrap" });

  // Setup config
  const config = new Config(logger);
  container.set("config", config);
  config.load({
    port: 3000,
    passwordSecret: "beeeeeeees",

    appName: "App",
    appBaseUrl: "localhost:3000",
    appEmail: "admin@app.localhost",
    appForceSsl: false,
    appLanguages: "en,fr",
    appDefaultLanguage: "en",
    appNotFoundTemplate: "errors/404",
    appErrorTemplate: "errors/500",

    mailImplementation: "console",

    authCookieKeys: "keyboardcat1,keyboardcat2",
    authLoginUrl: "/login",
    authUserProviderServiceKey: "services:user-provider",
    authCookieName: "authKey",
  });
  config.loadFromEnv();

  // Setup providers
  const providers = [
    require("./i18n"), //
    require("./app"),
    require("./errors"), //
    require("./template"), //
    require("./views"), //
    // require("./geo"), //
    // require("./db"),
    // require("./mail"), //
    // require("./jwt"),
    // require("./auth"),
  ].concat(userProviders);

  // For each provider call their pre-/post-/register methods in order
  logger.info("preRegister", { component: "bootstrap" });
  providers.forEach(function(provider) {
    if (provider.preRegister) provider.preRegister(container);
  });
  logger.info("register", { component: "bootstrap" });
  providers.forEach(function(provider) {
    if (provider.register) provider.register(container);
  });
  logger.info("postRegister", { component: "bootstrap" });
  providers.forEach(function(provider) {
    if (provider.postRegister) provider.postRegister(container);
  });
}
