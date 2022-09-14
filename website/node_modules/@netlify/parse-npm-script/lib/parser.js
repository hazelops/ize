const { parseCommand } = require("./parser-utils");

function getScriptName(command) {
  return command
    .replace(/^npm run\s/, "")
    .replace(/^npm\s/, "")
    .replace(/^yarn run\s/, "")
    .replace(/^yarn\s/, "")
    .replace(/(\w+)=('(.*)'|"(.*)"|(.*))\s/g, "");
}

function getSteps(pkg, scriptName) {
  const { scripts } = pkg;
  let current = [];

  if (scripts[scriptName]) {
    const rawCurrent = scripts[scriptName];
    const parsedCurrent = parseCommand(pkg, rawCurrent);
    // current = [scripts[scriptName]]
    current = [
      {
        name: scriptName,
        raw: rawCurrent,
        parsed: Array.isArray(parsedCurrent)
          ? flatten(parsedCurrent)
          : parsedCurrent,
      },
    ];
  }

  let preParsed = [];
  if (scripts[`pre${scriptName}`]) {
    // preParsed = [scripts[`pre${scriptName}`]]
    const rawPre = scripts[`pre${scriptName}`];
    const parsedPre = parseCommand(pkg, rawPre);
    preParsed = [
      {
        name: `pre${scriptName}`,
        raw: rawPre,
        parsed: Array.isArray(parsedPre) ? flatten(parsedPre) : parsedPre,
      },
    ];
  }
  // console.log('preParsed', preParsed)

  let postParsed = [];
  if (scripts[`post${scriptName}`]) {
    const rawPost = scripts[`post${scriptName}`];
    const parsedPost = parseCommand(pkg, rawPost);
    // postParsed = [scripts[`post${scriptName}`]]
    postParsed = [
      {
        name: `post${scriptName}`,
        raw: rawPost,
        parsed: Array.isArray(parsedPost) ? flatten(parsedPost) : parsedPost,
      },
    ];
  }
  // console.log('postParsed', postParsed)

  return preParsed.concat(current).concat(postParsed);
}

function parseNpmScript(pkg, command) {
  // console.log(`Processing command <${command}>`)
  const { scripts } = pkg;
  if (!scripts) {
    throw new Error(`No "scripts" field in package.json`);
  }
  const scriptName = getScriptName(command);

  if (!scripts[scriptName]) {
    throw new Error(`npm script "${scriptName}" not found`);
  }
  // console.log('scriptName', scriptName)
  const npmScriptSteps = getSteps(pkg, scriptName);
  // console.log('npm runs in this order', npmScriptSteps)

  const rawStepsInOrder = flatten(
    npmScriptSteps.map((x) => {
      return x.parsed;
    })
  );

  return {
    command: command,
    steps: npmScriptSteps,
    raw: rawStepsInOrder,
    combined: rawStepsInOrder.join(" && "),
  };
}

function flatten(arr) {
  return arr.reduce((flat, toFlatten) => {
    return flat.concat(
      Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten
    );
  }, []);
}

module.exports = parseNpmScript;
