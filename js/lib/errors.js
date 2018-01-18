class ApplicationError extends Error {}
exports.ApplicationError = ApplicationError;

class HttpError extends ApplicationError {
  constructor(statusCode, message = "Error") {
    super(message);
    this.statusCode = statusCode
    this.message = message;
  }
}
exports.HttpError = HttpError;

class BadRequestHttpError extends HttpError {
  constructor(message = "Bad Request") {
    super(400, message);
  }
}
exports.BadRequestHttpError = BadRequestHttpError;

class UnauthorizedHttpError extends HttpError {
  constructor(message = "Unauthorized") {
    super(401, message);
  }
}
exports.UnauthorizedHttpError = UnauthorizedHttpError;

class NotFoundHttpError extends HttpError {
  constructor(message = "Not Found") {
    super(404, message);
  }
}
exports.NotFoundHttpError = NotFoundHttpError;

class ServerErrorHttpError extends HttpError {
  constructor(message = "Server Error") {
    super(500, message);
  }
}
exports.ServerErrorHttpError = ServerErrorHttpError;
