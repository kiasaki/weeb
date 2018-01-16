class Config {
  constructor() {
    this.values = {};
  }

  get(key, default_) {
    return this.values[key] || default_;
  }

  set(key, value) {
    return (this.values[key] = value);
  }

  update(key, updateFn) {
    this.values[key] = updateFn(this.values[key]);
  }
}

module.exports = Config;
