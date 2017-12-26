'use strict';

var R = require('ramda');

// Sublasses need to implement
// - createEntityFromObject
// - createObjectFromEntity
function Repository(db) {
  this.db = db;
}
module.exports = Repository;

Repository.prototype.primaryKey = 'id';
Repository.prototype.table = null; // filled when subclassed
Repository.prototype.fields = []; // filled when subclassed

Repository.subclass = function(klass) {
  var parent = Repository;
  klass.prototype = Object.create(parent.prototype);
  klass.prototype.constructor = klass;
  klass.super_ = parent;
  R.forEach(function(key) {
    klass[key] = parent[key];
  }, R.keys(parent));
};

Repository.prototype.find = function(where, limit, offset) {
  var repo = this;
  var options = {};
  if (where) options.where = where;
  if (limit) options.limit = limit;
  if (offset) options.offset = offset;

  return repo.db.select(repo.table, options)
    .then(function(rows) {
      return R.map(function(obj) {
        return repo.createEntityFromObject(obj);
      }, rows);
    });
};

Repository.prototype.findOne = function(where) {
  return this.find(where, 1)
    .then(R.nth(0));
};

Repository.prototype.findById = function(pk) {
  var where = {};
  where[this.primaryKey] = pk;
  return this.findOne(where);
};

Repository.prototype.save = function(entity) {
  var pk = this.primaryKey;
  var obj = this.createObjectFromEntity(entity);
  var table = this.table;
  var fields = this.fields;

  if (obj[pk]) {
    var where = {}; where[pk] = obj[pk];
    return this.db.update(table, where, obj);
  } else {
    delete obj[pk];
    return this.db.insert(table, obj)
      .then(function(row) {
        // Update fields based on returned data
        R.forEach(function(key) {
          if (row[key]) {
            entity[key] = row[key];
          }
        }, fields);

        return entity;
      });
  }
};
