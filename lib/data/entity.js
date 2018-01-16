import { contains, ulidToUuid, uuidToUlid } = require("../utils");

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

Entity.EntityFields = {
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

module.exports = Entity;
