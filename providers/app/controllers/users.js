const {clone} = require("ramda");
const validator = require("../library/validator");
const User = require("../entities/user");

class UsersController {
    constructor(authService, userService) {
        this.authService = authService;
        this.userService = userService;

        this.edit = this.edit.bind(this);
        this.destroy = this.destroy.bind(this);
    }

    async edit(ctx) {
        let saved = false;
        let errors = [];
        const user = new User(clone(ctx.currentUser));

        if (ctx.method === "POST") {
            const body = ctx.request.body;
            user.name = body.name;
            user.email = body.email;

            errors = validator.validate([
                ["name", user.name, "required"],
                ["email", user.email, "required", "email", ["min", 6]],
                ["password", body.password, ["min", 8]],
            ]);

            // Now check for duplicate emails too (ignoring the currentUser)
            const duplicateUser = await this.userService.findByEmail(user.email);
            if (duplicateUser && duplicateUser.id != user.id) {
                errors.push("This email is already in use.");
            }

            // If password update was provided, hash password
            if (body.password) {
                user.password = await this.authService.hashPassword(body.password);
            }

            if (errors.length === 0) {
                // Save user updates
                await this.userService.update(user);
                // Make sure current views/render has the latest
                ctx.state.currentUser = user;
                saved = true;
            }
        }

        await ctx.render("users/edit", {user, errors, saved});
    }

    async destroy(ctx) {
        if (ctx.method === "POST") {
            await this.userService.destroy(ctx.currentUser.id);
            ctx.cookies.set("session_user_id", null, {overwrite: true});
            ctx.redirect("/");
            return;
        }

        await ctx.render("users/destroy");
    }
}

UsersController.dependencies = ["services:auth", "services:user"];

module.exports = UsersController;
