const bcrypt = require("bcrypt");

const bcryptSaltRounds = 13;

class AuthService {
    constructor(config) {
        this.config = config;
    }

    async hashPassword(password) {
        return await bcrypt.hash(password, bcryptSaltRounds);
    }

    async verifyPassword(password, passwordHash) {
        return await bcrypt.compare(password, passwordHash);
    }
}

AuthService.dependencyName = "services:auth";
AuthService.dependencies = ["config"];

module.exports = AuthService;
