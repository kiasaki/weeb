class MailConsoleService {
  send(payload) {
    console.log("<!EMAIL------------------------>");
    console.log(JSON.stringify(payload, null, 2));
    console.log("<!EMAIL------------------------>");
    return Promise.resolve();
  }
}

MailConsoleService.dependencyName = "services:mailConsole";
MailConsoleService.dependencies = [];
export default MailConsoleService;
