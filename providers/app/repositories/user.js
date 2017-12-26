const moment = require("moment");
const User = require("../entities/user");

class UserService {
    constructor(db) {
        this.db = db;
    }

    findById(id) {
        return this.db.findWhere(User, {id});
    }

    findByEmail(email) {
        return this.db.findWhere(User, {email});
    }

    create(entity) {
        return this.db.create(User, entity);
    }

    update(entity) {
        entity.updated = moment().utc().toDate();
        return this.db.update(User, entity);
    }

    destroy(entityId) {
        return this.db.knex.transaction(async function(trx) {
            const formIds = await trx("forms")
                .where("user_id", entityId)
                .pluck("id");

            // Delete all submissions
            await trx("submissions")
                .whereIn("form_id", formIds)
                .del();

            // Delete all forms
            await trx("forms")
                .where("user_id", entityId)
                .del();

            // Delete user
            await trx("users")
                .where("id", entityId)
                .del();
        });
    }
}

UserService.dependencyName = "services:user";
UserService.dependencies = ["db"];

module.exports = UserService;
