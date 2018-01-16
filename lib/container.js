class Container {
  constructor() {
    this.values = {};
  }

  get(key) {
    return this.values[key];
  }

  set(key, value) {
    return (this.values[key] = value);
  }
}

module.exports = Container;
