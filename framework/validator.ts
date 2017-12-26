import { forEach } from "ramda";

export function validate(spec) {
  const errors = [];

  forEach(property => {
    const name = property[0];
    const value = property[1];
    const validations = property.slice(2);
    let stop = false;

    forEach(validation => {
      // Limit ourselves to reporting 1 error per field
      if (stop) return;

      let validationName = "";
      let validationArgs = [];

      if (typeof validation === "string") {
        validationName = validation;
      } else {
        validationName = validation[0];
        validationArgs = validation.slice(1);
      }

      switch (validationName) {
        case "required":
          if (!value) {
            errors.push(`The ${name} is a required field.`);
            stop = true;
          }
          break;
        case "email":
          if (value && value.indexOf("@") === -1) {
            errors.push(`The ${name} doesn't look like a valid one.`);
            stop = true;
          }
          break;
        case "min":
          if (value && value.length < validationArgs[0]) {
            errors.push(
              `The ${name} must be at least ${validationArgs[0]} long.`,
            );
            stop = true;
          }
          break;
        default:
          throw new Error(
            "validatior: Unknown validation type: " + validationName,
          );
      }
    }, validations);
  }, spec);

  return errors;
}
