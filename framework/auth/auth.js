const thirtyDaysInMillis = 30 * 24 * 60 * 60 * 1000;

// Require a "service:user-provider" to be registered in container with:
// up.userIdFor(user)
// up.retriveById(id) -> returning a promise
// up.retriveByCredentials(credentials) -> returning a promise
// up.validateCredentials(user, credentials) -> returning a promise

class Auth {
  constructor(container, ctx) {
    this.container = container;
    this.ctx = ctx;
    this.user = null;

    const config = this.container.get("config");
    this.userProvider = this.container.get(
      config.get("authUserProviderServiceKey"),
    );
    this.cookieName = config.get("authCookieName");
  }

  isGuest() {
    return !this.user;
  }

  isLoggedIn() {
    return Boolean(this.user);
  }

  login(user, remember = false) {
    this.user = user;
    if (remember) {
      const id = this.userProvider.userIdFor(user);
      this.ctx.cookie.set(this.cookieName, id, {
        signed: true,
        expires: remember ? new Date(Date.now() + thirtyDaysInMillis) : 0,
      });
    }
  }

  attempt(credentials, login = false, remember = false) {
    let user;

    // Retrive user
    return this.userProvider
      .retriveByCredentials(credentials)
      .then(foundUser => {
        if (!foundUser) return false;

        // Validate user credentials
        return Promise.all([
          foundUser,
          userProvider.validateCredentials(user, credentials),
        ]);
      })
      .then(([user, valid]) => {
        if (!valid) return false;
        // Login if asked to do so
        if (login) this.login(user, remember);
        return true;
      });
  }

  validate(credentials) {
    return this.attempt(credentials, false, false);
  }

  logout() {
    this.ctx.cookies.set(this.cookieName, "", {
      signed: true,
      expires: new Date(Date.now() - 1),
    });
  }
}

export default Auth;
