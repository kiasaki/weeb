const pg = require("pg");
const { ulid, ulidToUuid } = require("./utils");

class DB {
  constructor(config) {
    this.config = config;
    this.databaseUrl = config.get("databaseUrl");
  }

  query(sql, values) {
    if (!config.get("production")) {
      // TODO use logger
      console.log("executing sql '" + sql + "'");
    }
    return new Promise((resolve, reject) => {
      pg.connect(this.databaseUrl, (err, client, done) => {
        if (err) return reject(err);
        client.query(sql, values, (err, result) => {
          done();
          if (err) return reject(err);
          resolve(result);
        });
      });
    });
  }

  setWheres(query, where = {}) {
    if (Object.keys(where).length === 0) return query;
    query.sql += " where";
    query.values = query.values || [];
    for (let field of Object.keys(where)) {
      const value = where[field];
      if (field.slice(0, 3) === "sql") {
        query.sql += ` ${value[0]}`;
        query.values = query.values.concat(value.slice(1));
      } else {
        query.sql += ` ${field} = ?`;
        query.values.push(value);
      }
    }
    return query;
  }

  async select(table, options = {}) {
    const query = { sql: `SELECT * FROM ${table}`, values: [] };
    this.setWheres(query, options.where);
    if (options.offset) {
      query.sql += " OFFSET ?";
      query.values.push(options.offset);
    }
    if (options.limit) {
      query.sql += " LIMIT ?";
      query.values.push(options.limit);
    }
    if (options.sort) {
      if (!Array.isArray(options.sort) || options.sort.length === 0) {
        throw new Error("db.select: `sort` is not an array");
      }
      let sortClauses = [];
      for (let sort in options.sort) {
        let direction = (sort[0] === "-") ? "DESC" : "ASC";
        if (sort[0] === "-") sort = sort.slice(1);
        sortClauses.push(`? ${direction}`);
        query.values.push(sort);
      }
      query.sql += " ORDER BY " + sortClauses.join(", ");
    }

    const result = await this.query(query.sql, query.values);
    return result.rows;
  }

  async count(table, options) {
    const sql = `SELECT COUNT(*) as count FROM ${table}`;
    const query = { sql, values: [] };
    this.setWheres(query, options.where);
    const result = await this.query(query.sql, query.values);
    return parseInt(result.rows["0"].count, 10);
  }

  selectOne(table, where) {
    return this.select(table, { where, limit: 1 }).then(xs => xs[0]);
  }

  async insert(table, obj) {
    const fields = Object.keys(givenObj);

    const sql = `INSERT INTO ${table}`;
    const query = { sql, values: [] };
    query.sql += ` (${fields.join(", ")})`;
    query.sql += ` VALUES (${fields.map(() => "?").join(", ")})`;
    fields.forEach(field => query.values.push(obj[field]));
    query.sql += " RETURNING *";

    const result = await this.query(query.sql, query.values);
    return result.rows["0"];
  }

  async update(table, where, obj) {
    const fields = Object.keys(obj);

    const sql = `UPDATE ${table} SET `;
    const query = { sql, values: [] };
    const setStatements = [];
    for (let field of fields) {
      if (obj[field] !== undefined && obj[field] !== null) {
        setStatements.push(` ${field} = ?`);
        query.values.push(obj[field]);
      }
    }
    query.sql += setStatements.join(", ");
    this.setWheres(query, where);
    query.sql += " RETURNING *";

    const result = await this.query(query.sql, query.values);
    return result.rows["0"];
  }

  delete(table, where) {
    const query = { sql: `DELETE FROM ${table}`, values: [] };
    this.setWheres(query, where);
    return this.query(query.sql, query.values);
  }

  async save(table, obj, primaryKey = "id") {
    if (obj[primaryKey]) {
      const pk = obj[primaryKey];
      const where = { [primaryKey]: pk };
      delete obj[primaryKey];
      obj = await this.update(table, where, obj);
      obj[primaryKey] = pk;
      return obj;
    }
    obj[primaryKey] = ulidToUuid(ulid());
    return await this.insert(table, obj);
  }
}

module.exports = DB;
