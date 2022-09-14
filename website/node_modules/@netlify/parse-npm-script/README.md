# Parse NPM Script

Parse a given `npm` script command from `package.json` & return information.

- What the `script` lifecycle looks like
- How the command **really** resolves
- & What the final command run is

Useful for parsing dependencies actually used and figuring out wtf is happening in a complex `package.json`.
It supports extracting values from npm-run-all.

## Usage

```js
const path = require("path");
const util = require("util");
const parse = require("@netlify/parse-npm-script");
const { readFileSync } = require("fs-extra");

/* your package.json file */
const packageJSON = readFileSync(
  path.join(__dirname, "tests/fixtures/one.json")
);

async function runParser() {
  const parsed = await parse(packageJSON, "npm run build");
  console.log(
    util.inspect(parsed, {
      showHidden: false,
      depth: null,
    })
  );
}

runParser();

/* Parsed contents
{
  command: 'npm run build',
  steps: [{
      name: 'prebuild',
      raw: 'echo a && npm run foo',
      parsed: ['echo a', 'echo foo']
    },
    {
      name: 'build',
      raw: 'echo b && npm run cleanup',
      parsed: ['echo b', 'echo cleanup']
    },
    {
      name: 'postbuild',
      raw: 'echo c',
      parsed: 'echo c'
    }
  ],
  raw: ['echo a', 'echo foo', 'echo b', 'echo cleanup', 'echo c'],
  combined: 'echo a && echo foo && echo b && echo cleanup && echo c'
}
*/
```

## Example:

Parsing a `package.json`

```json
{
  "name": "parse-npm-script",
  "scripts": {
    "foo": "echo foo",
    "cleanup": "echo cleanup",
    "prebuild": "echo a && npm run foo",
    "build": "echo b && npm run cleanup",
    "postbuild": "echo c"
  },
  "author": "David Wells",
  "license": "MIT"
}
```

Will result in this output:

```js
{
  command: 'npm run build',
  steps: [{
      name: 'prebuild',
      raw: 'echo a && npm run foo',
      parsed: ['echo a', 'echo foo']
    },
    {
      name: 'build',
      raw: 'echo b && npm run cleanup',
      parsed: ['echo b', 'echo cleanup']
    },
    {
      name: 'postbuild',
      raw: 'echo c',
      parsed: 'echo c'
    }
  ],
  raw: ['echo a', 'echo foo', 'echo b', 'echo cleanup', 'echo c'],
  combined: 'echo a && echo foo && echo b && echo cleanup && echo c'
}
```
