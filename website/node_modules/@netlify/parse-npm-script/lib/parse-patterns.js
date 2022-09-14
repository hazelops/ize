/**
 * @module index
 * @author Toru Nagashima
 * @copyright 2015 Toru Nagashima. All rights reserved.
 * See LICENSE file in root directory for full license.
 */
"use strict";

//------------------------------------------------------------------------------
// Requirements
//------------------------------------------------------------------------------

const shellQuote = require("shell-quote");
const matchTasks = require("./match-tasks");

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

const ARGS_PATTERN = /\{(!)?([*@]|\d+)([^}]+)?}/g;

/**
 * Converts a given value to an array.
 *
 * @param {string|string[]|null|undefined} x - A value to convert.
 * @returns {string[]} An array.
 */
function toArray(x) {
  if (x == null) {
    return [];
  }
  return Array.isArray(x) ? x : [x];
}

/**
 * Replaces argument placeholders (such as `{1}`) by arguments.
 *
 * @param {string[]} patterns - Patterns to replace.
 * @param {string[]} args - Arguments to replace.
 * @returns {string[]} replaced
 */
function applyArguments(patterns, args) {
  const defaults = Object.create(null);

  return patterns.map((pattern) =>
    pattern.replace(ARGS_PATTERN, (whole, indirectionMark, id, options) => {
      if (indirectionMark != null) {
        throw Error(`Invalid Placeholder: ${whole}`);
      }
      if (id === "@") {
        return shellQuote.quote(args);
      }
      if (id === "*") {
        return shellQuote.quote([args.join(" ")]);
      }

      const position = parseInt(id, 10);
      if (position >= 1 && position <= args.length) {
        return shellQuote.quote([args[position - 1]]);
      }

      // Address default values
      if (options != null) {
        const prefix = options.slice(0, 2);

        if (prefix === ":=") {
          defaults[id] = shellQuote.quote([options.slice(2)]);
          return defaults[id];
        }
        if (prefix === ":-") {
          return shellQuote.quote([options.slice(2)]);
        }

        throw Error(`Invalid Placeholder: ${whole}`);
      }
      if (defaults[id] != null) {
        return defaults[id];
      }

      return "";
    })
  );
}

/**
 * Parse patterns.
 * In parsing process, it replaces argument placeholders (such as `{1}`) by arguments.
 *
 * @param {string|string[]} patternOrPatterns - Patterns to run.
 *      A pattern is a npm-script name or a Glob-like pattern.
 * @param {string[]} args - Arguments to replace placeholders.
 * @returns {string[]} Parsed patterns.
 */
function parsePatterns(patternOrPatterns, args) {
  const patterns = toArray(patternOrPatterns);
  const hasPlaceholder = patterns.some((pattern) => ARGS_PATTERN.test(pattern));

  return hasPlaceholder ? applyArguments(patterns, args) : patterns;
}

module.exports = parsePatterns;
