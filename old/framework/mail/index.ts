export function register(container) {
  container.load(require("./mail-console"));
  container.load(require("./mail-smtp"));
  container.load(require("./mail-postmark"));
  container.load(require("./mail"));
}
