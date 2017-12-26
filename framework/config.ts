import * as fs from "fs";

import Logger from "./logger";
import { camelCase } from "./case";

class Config {
  constructor(logger = new Logger()) {
    this.config = {};
    this.logger = logger.withSubContext({
      component: "config",
    });
  }

  get(key, default_) {
    return this.config[key] || default_;
  }

  set(key, value) {
    this.config[key] = value;
  }

  unset(key) {
    delete this.config[key];
  }

  reset() {
    this.config = {};
  }

  load(values) {
    this.config = Object.assign(this.config, values);
  }

  loadFromEnv() {
    this.logger.info("loading from env");
    Object.keys(process.env).forEach(envKey => {
      this.set(camelCase(envKey), process.env[envKey]);
    });
  }

  loadFromFile(filename) {
    this.logger.info("loading from file", { filename });
    try {
      const contents = fs.readFileSync(filename, {
        encoding: "utf8",
      });
      const values = JSON.parse(contents);
      this.load(values);
      return true;
    } catch (e) {
      return false;
    }
  }
}

export default Config;
