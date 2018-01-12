const { Application } = require("./lib.js");

const app = new Application();
app.addApplication("app", __dirname + "/app");
app.router.get("/", (req, res) => res.template(200, "home", { name: "Joe" }));
app.start();
