const parse = require("./parser");

/**
 * Parse npm package script command
 * @param  {object} packageJson parsed package.json object
 * @param  {string}  command - npm script command to resolve & parse
 * @return {Promise} - resolved command data
 */
function parseNpmScript(packageJson, command) {
  return parse(packageJson, command);
}

module.exports = parseNpmScript;
