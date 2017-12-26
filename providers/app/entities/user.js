const Entity = require("../library/entity");

class User extends Entity {
}

User.table = "users";
User.prototype.defaults = {
    created: Entity.newDate,
    updated: Entity.newDate,
};
User.prototype.fields = [
    "id",
    "name",
    "email",
    "password",
    "created",
    "updated",
];
User.prototype.privateFields = [
    "password",
];

module.exports = User;
