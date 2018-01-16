const append = a => b => b.concat([a]);
const contains = item => list => list.indexOf(item) !== -1;

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

const paramCase = value => noCase(value, "-");
const snakeCase = value => noCase(value, "_");
const camelCase = value =>
  noCase(value).replace(/ (.)/g, (_, v) => v.toUpperCase());

const caseKeys = (caseFn) => (obj) => {
  const newObj = {};
  for (let key of Object.keys(obj)) {
    newObj[caseFn(key)] = obj[key];
  }
  return newObj;
};
const snakeCaseKeys = caseKeys(snakeCase);
const camelCaseKeys = caseKeys(camelCase);

function hash(str) {
  let hash = 5381;
  let i = str.length;
  while(i) {
    hash = (hash * 33) ^ str.charCodeAt(--i);
  }

  // Make positive using an unsigned bitshift of 0
  return hash >>> 0;
}

const ULID_CHARS = "0123456789ABCDEFGHJKMNPQRSTVWXYZ";
const ULID_CHARS_LENGTH = ULID_CHARS.length;

function ulid() {
  let str = "";
  let now = Date.now();
  let mod;
  for (let i = 10; i > 0; i--) {
    mod = now % ULID_CHARS_LENGTH;
    str = ULID_CHARS.charAt(mod) + str;
    now = (now - mod) / ULID_CHARS_LENGTH;
  }
  for (let i = 16; i > 0; i--) {
    let rand = Math.floor(Math.random() * ULID_CHARS_LENGTH);
    if (rand >= ULID_CHARS_LENGTH) rand = ULID_CHARS_LENGTH-1;
    str = str + ULID_CHARS.charAt(rand);
  }
  return str;
}

function ulidToUuid(id) {
  return [
    id.slice(0, 8),
    id.slice(8, 12),
    id.slice(12, 16),
    id.slice(16, 20),
    id.slice(20, 26) + "000000",
  ].join("-");
}

function uuidToUlid(id) {
  return id.replace("-", "").slice(0, 26);
}

exports.append = append;
exports.contains = contains;

exports.noCase = noCase;
exports.paramCase = paramCase;
exports.snakeCase = snakeCase;
exports.camelCase = camelCase;
exports.snakeCaseKeys = snakeCaseKeys;
exports.camelCaseKeys = camelCaseKeys;

exports.hash = hash;
exports.ulid = ulid;
exports.ulidToUuid = ulidToUuid;
exports.uuidToUlid = uuidToUlid;
