function noCase(str, replacement = " ") {
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

const snakeCase = value => noCase(value, "_");
const paramCase = value => noCase(value, "-");
const camelCase = value =>
  noCase(value).replace(/ (.)/g, (_, v) => v.toUpperCase());

function hash(str) {
  let hash = 5381;
  let i = str.length;
  while(i) {
    hash = (hash * 33) ^ str.charCodeAt(--i);
  }

  // Make positive using an unsigned bitshift of 0
  return hash >>> 0;
}

exports.noCase = noCase;
exports.snakeCase = snakeCase;
exports.paramCase = paramCase;
exports.camelCase = camelCase;
exports.hash = hash;
