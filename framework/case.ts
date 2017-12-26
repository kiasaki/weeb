export function noCase(str, replacement = " ") {
  if (!str) return "";

  const replace = (match, index, value) => {
    if (index === 0 || index === value.length - match.length) {
      return "";
    }
    return replacement;
  };

  return String(str)
    .replace(/([a-z])([A-Z0-9])/g, "$1 $2")
    .replace(/([A-Z]+)([A-Z0-9])/g, "$1 $2")
    .replace(/[^a-zA-Z0-9]+/g, replace)
    .toLowerCase();
}

export const snakeCase = value => noCase(value, "_");
export const paramCase = value => noCase(value, "-");
export const camelCase = value =>
  noCase(value).replace(/ (.)/g, (_, v) => v.toUpperCase());
