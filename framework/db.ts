const R = require("ramda");
const pg = require("pg");
const uuid = require("node-uuid");
const humps = require("humps");
const moment = require("moment");
const squel = require("squel").useFlavour("postgres");
const debug = require("debug")("lotad:db");

squel.registerValueHandler(Date, date => {
  const formattedDate = moment(date)
    .utc()
    .format("YYYY-MM-DD hh:mm:ss");
  return `"${formattedDate}"`;
});

squel.registerValueHandler("undefined", () => {
  return null;
});

squel.registerValueHandler("object", jsonObj => {
  return JSON.stringify(jsonObj);
});

function DB(config) {
  this.dbUrl = config.get("database_url");
}

DB.dependencyName = "db";
DB.dependencies = ["config"];

module.exports = DB;

DB.prototype.query = function(text, values) {
  if (process.env.NODE_ENV !== "production") {
    debug(`executing sql "${text}"`);
  }
  const dbUrl = this.dbUrl;

  return new Promise((resolve, reject) => {
    // Aquire a connection
    pg.connect(dbUrl, (err, client, done) => {
      if (err) {
        return reject(err);
      }

      // Send out query
      client.query(text, values, (err, result) => {
        // Release connection back to pool
        done();
        if (err) {
          return reject(err);
        }

        resolve(result);
      });
    });
  });
};

DB.prototype.setWheres = function(sql, wheres) {
  if (!wheres) {
    return sql;
  }

  R.forEach(key => {
    if (key.slice(0, 3) === "sql") {
      sql = sql.where(...wheres[key]);
    } else {
      sql = sql.where(`${humps.decamelize(key)} = ?`, wheres[key]);
    }
  }, R.keys(wheres));

  return sql;
};

DB.prototype.select = function(table, options) {
  let sql = squel.select().from(table);

  if (options && options.where) {
    sql = this.setWheres(sql, options.where);
  }
  if (options && options.offset) {
    sql = sql.offset(options.offset);
  }
  if (options && options.limit) {
    sql = sql.limit(options.limit);
  }
  if (options && options.sort) {
    if (!Array.isArray(options.sort)) {
      return Promise.reject(
        new Error("sort passed to db.selectOne must be an array")
      );
    }
    R.forEach(column => {
      let asc = true;
      if (column[0] === "-") {
        asc = false;
        column = column.slice(1);
      }
      sql = sql.order(column, asc);
    });
  }

  const paramQuery = sql.toParam();
  return this.query(paramQuery.text, paramQuery.values)
    .then(R.prop("rows"))
    .then(R.map(humps.camelizeKeys));
};

DB.prototype.count = function(table, options) {
  let sql = squel
    .select()
    .from(table)
    .field("COUNT(*) as count");
  sql = this.setWheres(sql, options.where);

  const paramQuery = sql.toParam();
  return this.query(paramQuery.text, paramQuery.values)
    .then(R.path(["rows", "0", "count"]))
    .then(parseInt);
};

DB.prototype.selectOne = function(table, where) {
  return this.select(table, { where, limit: 1 })
    .then(R.nth(0))
    .then(humps.camelizeKeys);
};

DB.prototype.executeSquelAndReturnFirst = function(sql) {
  const paramQuery = sql.toParam();
  return this.query(
    paramQuery.text,
    R.map(v => (v === "null" ? null : v), paramQuery.values)
  ).then(R.pipe(R.prop("rows"), R.nth(0), humps.camelizeKeys));
};

DB.prototype.insert = function(table, what) {
  const sql = squel
    .insert()
    .into(table)
    .setFieldsRows([undefToNull(humps.decamelizeKeys(what))])
    .returning("*");

  return this.executeSquelAndReturnFirst(sql);
};

DB.prototype.update = function(table, where, what) {
  const params = R.fromPairs(
    R.filter(pair => {
      return pair[1] === false || pair[1];
    }, R.toPairs(what))
  );
  let sql = squel
    .update()
    .table(table)
    .setFields(undefToNull(humps.decamelizeKeys(params)));

  sql = this.setWheres(sql, where);

  return this.executeSquelAndReturnFirst(sql);
};

DB.prototype.delete = function(table, where) {
  let sql = squel.delete().from(table);

  sql = this.setWheres(sql, where);

  const paramQuery = sql.toParam();
  return this.query(paramQuery.text, paramQuery.values);
};

DB.prototype.save = function(table, entity) {
  if (entity.id) {
    const id = entity.id;
    delete entity.id;
    return this.update(table, { id }, entity);
  }
  entity.id = uuid.v4();
  return this.insert(table, entity).then(R.always(entity));
};

function undefToNull(obj) {
  return R.mapObjIndexed(value => {
    return typeof value === "undefined" ? null : value;
  }, obj);
}

exports.register = function(container) {
  var dbUrl = container.get("config").get("database_url");

  var db = new DB(dbUrl);

  container.set("db", db);
};
