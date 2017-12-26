import * as R from "ramda";

import Logger from "./logger";

declare interface ConstructorWithDependencies {
  new (...a: Array<any>): any;
  dependencies: Array<string>;
}

declare interface ConstructorWithDependencyName {
  new (...a: Array<any>): any;
  dependencyName: string;
  dependencies: Array<string>;
}

// Creates an instance of a class, passing args to the constructor
// ES6 constructors need the `new` keyword to be used, so, no .apply
// magic can be done here until WebReflection lands in Node.js stable.
function applyToConstructor(
  constructor: (...a: Array<any>) => any,
  args: Array<any>,
) {
  return R.apply(R.construct(constructor), args);
}

var nextIdCounter = 1;
function nextId() {
  return nextIdCounter++;
}

class Container {
  id: number;
  contents: { [key: string]: any };
  logger: Logger;

  constructor() {
    this.id = nextId();
    this.contents = {};
    this.setLogger(new Logger());
  }

  setLogger(logger: Logger): void {
    this.logger = logger.withSubContext({
      component: "container",
    });
  }

  get = (name: string): any => {
    if (!(name in this.contents)) {
      throw Error(
        "Container #" + this.id + " has nothing registered for key " + name,
      );
    }
    return this.contents[name];
  };

  set = (name: string, instance: any) => {
    this.logger.info("set instance", { name });
    this.contents[name] = instance;
    return instance;
  };

  create(klass: ConstructorWithDependencies) {
    if (!("length" in klass.dependencies)) {
      throw new Error(
        "Invariant: container can't resolve a class without a dependencies",
      );
    }

    const dependencies: Array<string> = [];
    R.forEach((dependencyName: string) => {
      dependencies.push(this.get(dependencyName));
    }, klass.dependencies);

    return applyToConstructor(klass, dependencies);
  }

  load(klass: ConstructorWithDependencyName) {
    if (typeof klass.dependencyName !== "string") {
      throw new Error(
        "Invariant: container can't resolve a class without a name",
      );
    }

    const instance = this.create(klass);
    return this.set(klass.dependencyName, instance);
  }

  unset(name: string) {
    delete this.contents[name];
  }

  reset() {
    this.contents = {};
  }
}

export default Container;
