import * as Koa from "koa";

import Config from "./config";
import Container from "./container";

class I18n {
  static dependencyName = "i18n";
  static dependencies = ["config"];

  config: Config;
  translations: { [language: string]: { [key: string]: string } };

  constructor(config: Config) {
    this.config = config;
    this.translations = {};
  }

  addTranslations(language: string, translations: object) {
    this.translations[language] = Object.assign(
      this.translations[language],
      flattenKeys(translations),
    );
  }

  t(lang: string, key: string) {
    return this.translations[lang][key];
  }
}

export default I18n;

export function register(container: Container) {
  const app = container.get("app");
  const config = container.get("config");
  const i18n = container.load(I18n);

  app.addTranslationsFolder = (folder: string) => {
    const languages = config.get("appLanguages", "en").split(",");
    languages.forEach((language: string) => {
      i18n.addTranslations(
        language,
        require(path.join(folder, language + ".json")),
      );
    });
  };

  container.get("app").use(async function(ctx: Koa.Context, next: Function) {
    ctx.t = (key: string) =>
      i18n.t(ctx.state.lang || config.get("appDefaultLanguage"), key);
    await next();
  });
}

function flattenKeys(obj, prefix = "") {
  let flattened = {};
  const newPrefix = prefix + (prefix ? "." : "");
  Object.keys(obj).forEach(key => {
    if (obj[key] instanceof Object) {
      flattened = Object.assign(
        {},
        flattened,
        flattenKeys(obj[key], newPrefix + key),
      );
    } else {
      flattened[newPrefix + key] = obj[key];
    }
  });
  return flattened;
}
