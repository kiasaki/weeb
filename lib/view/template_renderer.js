const FileSystem = require("../file_system");
const { hash } = require("../util");

class TemplateRenderer {
  constructor(config) {
    this.config = config;
    this.templateCache = {};
    const cacheFiles = config.get("production");
    this.fs = new FileSystem(cacheFiles);
  }

  compile(templateSource) {
    const code = "var p=[],print=function(){p.push.apply(p,arguments);};" +
      "with(context||{}){p.push('" +
      templateSource
        .replace(/[\r\t\n]/g, " ")
        .split("<%").join("\t")
        .replace(/((^|%>)[^\t]*)'/g, "$1\r")
        .replace(/\t=(.*?)%>/g, "',$1,'")
        .split("\t").join("');")
        .split("%>").join("p.push('")
        .split("\r").join("\\'")
      + "');}return p.join('');";
    return new Function("context", code);
  }

  async render(name, context) {
    const directories = this.config.get("templateFolders");
    const path = await this.fs.findInDirectories(name + ".html", directories);
    if (!path) {
      throw new Error("Template named '" + name + "' not found");
    }
    const templateSource = await this.fs.read(path);
    const templateSourceHash = hash(templateSource);
    const template = this.templateCache[templateSourceHash] || this.compile(templateSource);
    if (this.config.get("production")) {
      this.templateCache[templateSourceHash] = template;
    }
    return template(context);
  }
}

module.exports = TemplateRenderer;
