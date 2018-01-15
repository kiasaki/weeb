import R from "ramda";
import fs from "fs";
import path from "path";

class GeoService {
  constructor() {
    const countriesPath = path.join(__dirname, "data", "countries.txt");
    const statesPath = path.join(__dirname, "data", "states.txt");
    this._countries = R.map(
      countryLine => {
        return {
          code: countryLine.split(":")[0],
          name: countryLine.split(":")[1],
        };
      },
      fs
        .readFileSync(countriesPath, { encoding: "utf8" })
        .trim()
        .split("\n"),
    );
    this._states = {};
    R.forEach(
      stateLine => {
        const stateLineParts = stateLine.split(":");
        this._states[stateLineParts[0]] = this._states[stateLineParts[0]] || [];
        this._states[stateLineParts[0]].push({
          code: stateLineParts[1],
          name: stateLineParts[2],
          countryId: stateLineParts[0],
        });
      },
      fs
        .readFileSync(statesPath, { encoding: "utf8" })
        .trim()
        .split("\n"),
    );
  }

  countries() {
    return this._countries;
  }

  country(code) {
    const country = R.find(R.propEq("code", code), this._countries);
    if (!country) {
      throw new Error(`Country for code "${code}" not found`);
    }
    return country;
  }

  statesForCountry(code) {
    return code in this._states ? this._states[code] : [];
  }

  state(countryCode, stateCode) {
    const states = this.statesForCountry(countryCode);
    if (states.length === 0) {
      return "";
    }

    const state = R.find(R.propEq("code", stateCode), states);

    if (!state) {
      throw new Error(`State for code "${countryCode}" not found`);
    }
    return state;
  }
}
GeoService.dependencyName = "services:geo";
GeoService.dependencies = [];
export default GeoService;
