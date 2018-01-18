const R = require('ramda');

class Entity {
    constructor(params, skipDefaults) {
        const defaults = skipDefaults ? {} : (this.defaults || {});
        R.mapObjIndexed((v, k) => {
            this[k] = (typeof v === 'function') ? v() : v;
        }, R.merge(defaults, params));
    }

    toJson(includePrivate) {
        const privFields = this.privateFields || [];
        const pred = R.compose(R.not, R.curry(R.contains)(R.__, privFields));
        const fields = R.filter(pred, this.fields);

        return R.pickAll((includePrivate ? this.fields : fields), this);
    }

    toObject() {
        return this.toJson(true);
    }
}

Entity.newDate = function() {
    return new Date();
};

module.exports = Entity;
