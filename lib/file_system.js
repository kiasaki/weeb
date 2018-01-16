const fs = require("fs");
const path = require("path");
const { promisify } = require("util");

class FileSystem {
  constructor(cached = false) {
    this.cached = cached;
    this.fileCache = {};
  }

  async findInDirectories(name, folders) {
    for (let folder of folders) {
      try {
        const fullPath = path.join(folder, name);
        await promisify(fs.stat)(fullPath);
        return fullPath;
      } catch (err) {
        if (err.code === "ENOENT") {
          continue;
        }
        throw err;
      }
    }
    return null;
  }

  async isDirectory(name) {
    const stat = await promisify(fs.stat)(name);
    return stat.isDirectory();
  }

  async read(name) {
    if (name in this.fileCache) return this.fileCache[name];

    const contents = await promisify(fs.readFile)(name, { encoding: "utf8" });
    this.fileCache[name] = contents;
    return contents;
  }
}

module.exports = FileSystem;
