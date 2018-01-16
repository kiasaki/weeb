const { Project } = require("./lib");

const app = new Project();
app.addApplication("app", __dirname + "/app");
app.router.get("/", (req, res) => res.template(200, "home", { name: "Joe" }));
app.start();
