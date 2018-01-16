class Repository {
  constructor(db) {
    this.db = db;

    this.primaryKey = "id";
    this.table = null;
    this.entityClass = null;
  }

  async find(where, limit, offset) {
    if (!this.entityClass) throw new Error("Entity is missing entityClass");
    const options = {};
    if (where) options.where = where;
    if (limit) options.limit = limit;
    if (offset) options.offset = offset;

    const rows = await this.db.select(this.table, options);

    return rows.map(this.entityClass.fromDatabase.bind(this.entityClass));
  }

  findOne(where) {
    return this.find(where, 1).then(xs => xs[0]);
  }

  findById(pk) {
    return this.findOne({ [this.primaryKey]: pk });
  }

  save(entity) {
    const obj = entity.toDatabase();
    return this.entityClass.fromDatabase(
      this.db.save(table, obj, this.primaryKey)
    );
  }
}

module.exports = Repository;
