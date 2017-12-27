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
const { bootstrap, Container } = require('weeb');
const container = new Container();
bootstrap(container, [
  "./app"
]);

const app = container.get("app");

app.start();
```

### license

MIT. See `LICENSE` file.
