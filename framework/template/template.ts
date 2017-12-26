const R = require("ramda");
const fs = require("fs");
const path = require("path");
const hogan = require("hogan.js");

class TemplateService {
  constructor(config) {
    this.config = config;
    this.templates = {};
    this.partials = {};
    this.globals = {};
  }

  render(templateName, locals) {
    if (!(templateName in this.templates)) {
      throw new Error(`Can't find template ${templateName}`);
    }
    const template = this.templates[templateName];
    return template.render(R.merge(this.globals, locals), this.partials);
  }

  fetchTemplate(templateName) {
    const templatePath = path.join(this.config.get("root"), templateName);
    const rawTemplate = fs.readFileSync(templatePath, {
      encoding: "utf8",
    });
    return hogan.compile(rawTemplate);
  }

  loadTemplate(templateName, templatePath) {
    this.templates[templateName] = this.fetchTemplate(templatePath);
  }

  loadPartial(partialName, templatePath) {
    this.partials[partialName] = this.fetchTemplate(templatePath);
  }

  setGlobal(key, value) {
    this.globals[key] = value;
  }
}

TemplateService.dependencyName = "services:template";
TemplateService.dependencies = ["config"];
module.exports = TemplateService;
