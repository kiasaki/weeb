# weeb

_A weeb framework. For quickly building ambitious weeb applications._

### intro

### features

**there**

- Authentication
- Controllers
- DI / IOC Container
- Database CRUD
- Database Querying
- Encryption
- I18n
- Logging
- Mails
- Middewares
- Pluggable Providers
- Routing
- Templates
- Validation

**upcomming**

- Background Jobs
- Cache
- Cron Jobs
- Database Migrations
- Deployment
- File Storage
- Security
- Sessions
- Subscriptions Billing
- Test Helpers

### usage

### concepts

Everything is in the DI container, everything is a provider, everything
is optional / pluggable.

### starting a new project

```
npm init --yes
npm install --save weeb
```

Create an `index.js` file with:

```js
const { Application } = require("weeb");

const app = new Application();
app.addApplication("app", __dirname + "/app");
app.start();
```

Start the web server by running:

```
node index.js
```

### license

MIT. See `LICENSE` file.
