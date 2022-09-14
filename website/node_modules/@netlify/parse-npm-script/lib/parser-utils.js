const matchTasks = require("./match-tasks");
const parseCliArgs = require("./parse-cli-args");
const parsePatterns = require("./parse-patterns");

const __DEV__ = false;
const noOp = () => {};
let LOGGER = noOp;
if (__DEV__) {
  LOGGER = (padLevel, ...args) => {
    console.log(`${"  ".repeat(padLevel)}${args.join(" ")}`);
  };
}

const COMPLEX_BASH_COMMANDS = [
  // complex command, do not parse further
  "bash",
  "sh",
  "ssh",
  ".",
  "source",
  "su",
  "sudo",
  "cd", // TODO: try parse simple cd, to unwrap further
  "if",
  "eval",
  "cross-env",
];

const SIMPLE_BASH_COMMAND = [
  "node",
  "npx",
  "git",
  "rm",
  "mkdir",
  "echo",
  "cat",
  "exit",
  "kill",
];

const PACKAGE_INSTALLER = ["npm", "yarn"];
const NPM_RUN_ALL = ["npm-run-all", "run-p", "run-s"];

const REGEXP_ESCAPE = /\\/g;
const REGEXP_QUOTE = /[" ]/g;
const wrapJoinBashArgs = (args) =>
  args
    .map((arg) => {
      return `"${arg
        .replace(REGEXP_ESCAPE, "\\\\")
        .replace(REGEXP_QUOTE, "\\$&")}"`;
    })
    .join(" ");

const warpBashSubShell = (command) => `(
${indentLine(command, "  ")}
)`;

const ENV_REGEX = /(\w+)=('(.*)'|"(.*)"|(.*))/;

function parseCommand(pkg, scriptString, level, log = LOGGER) {
  log(level, "[parseCommand]", `input: <${scriptString}>`);
  scriptString = scriptString.trim();

  const [leadingCommand, secondCommand, ...additionalCommands] = scriptString
    .split(" ")
    .filter((cmd) => !cmd.match(ENV_REGEX) && cmd !== "cross-env");

  if (NPM_RUN_ALL.includes(leadingCommand)) {
    log(level, "npm-run-all command");

    const args = parseCliArgs(
      [secondCommand, ...additionalCommands],
      {},
      {
        packageConfig: { test: pkg.scripts },
        singleMode: leadingCommand !== "npm-run-all",
      }
    );

    const commands = args.groups.reduce(
      (prev, { patterns }) => [...prev, ...patterns],
      []
    );

    const tasks = matchTasks(
      Object.keys(pkg.scripts),
      parsePatterns(commands, [])
    );

    return tasks.map((command) => {
      return parseCommand(pkg, `npm run ${command}`, level + 1, log) || command;
    });
  }

  // Check for complex bash
  if (
    COMPLEX_BASH_COMMANDS.includes(leadingCommand) ||
    leadingCommand.startsWith("./")
  ) {
    log(level, "✓ directly executable complex command, return");
    return scriptString;
  } else {
    log(level, `? not directly executable complex command: ${leadingCommand}`);
  }

  // Check for combo commands
  if (scriptString.includes(" && ")) {
    log(level, "✓ combo command, split");

    const subCommandList = scriptString.split(" && ");
    return subCommandList.map((command) => {
      return parseCommand(pkg, command, level + 1, log) || command;
    });
    /*
    return warpBashSubShell(subCommandList.map((command) => {
      return parseCommand(pkg, command, level + 1, log) || command
    }).join('\n'))
    */
  } else {
    log(level, `? not combo command, I guess`);
  }

  if (SIMPLE_BASH_COMMAND.includes(leadingCommand)) {
    log(level, "✓ directly executable simple command, return");

    return scriptString;
  } else {
    log(level, `? not directly executable simple command: ${leadingCommand}`);
  }

  // TODO: consider allow package dependency command

  if (PACKAGE_INSTALLER.includes(leadingCommand)) {
    if (secondCommand === "run") {
      log(level, "✓ package script, parse");
      const [scriptName, ...extraArgs] = additionalCommands;
      extraArgs[0] === "--" && extraArgs.shift();
      return parsePackageScript(
        pkg,
        scriptName,
        extraArgs.join(" "),
        level + 1,
        log
      );
    }

    if (secondCommand === "test" || secondCommand === "t") {
      log(level, "✓ package test script, parse");
      const [...extraArgs] = additionalCommands;
      extraArgs[0] === "--" && extraArgs.shift();
      return parsePackageScript(
        pkg,
        "test",
        extraArgs.join(" "),
        level + 1,
        log
      );
    }
    if (secondCommand === "start") {
      log(level, "✓ package test script, parse");
      const [...extraArgs] = additionalCommands;
      extraArgs[0] === "--" && extraArgs.shift();
      return parsePackageScript(
        pkg,
        "start",
        extraArgs.join(" "),
        level + 1,
        log
      );
    }
    if (leadingCommand === "yarn") {
      log(level, "✓ yarn package script, parse");

      if (secondCommand === "workspace") {
        log(level, "✓ yarn workspace package script, bail");
        return scriptString;
      }
      return parsePackageScript(
        pkg,
        secondCommand,
        additionalCommands.join(" "),
        level + 1,
        log
      );
    }
  } else {
    log(level, "? unknown npm/yarn script");
  }

  log(level, "? unknown script, bail");
  return scriptString;
}

function parsePackageScript(
  pkg,
  scriptName,
  extraArgsString = "",
  level,
  log = LOGGER
) {
  log(
    level,
    "[parsePackageScript]",
    `script name: <${scriptName}>, extra: ${extraArgsString}`
  );

  if (scriptName === "-w" || extraArgsString.split(" ").includes("-w")) {
    log(level, "? workspace script. bailing.");
    return `${scriptName} ${extraArgsString}`;
  }

  const scriptString = pkg.scripts[scriptName];
  if (!scriptString) {
    throw new Error(
      `[parsePackageScript] missing script with name: ${scriptName}`
    );
  }

  const otherScriptString = [scriptString, extraArgsString]
    .filter(Boolean)
    .join(" ");
  const command = parseCommand(pkg, otherScriptString, level + 1, log);
  if (command) {
    return command;
  }

  log(level, "? unexpected script, bail to npm run");

  return [`${otherScriptString}`, extraArgsString].filter(Boolean).join(" -- ");
  // return [`npm run ${scriptName}`, extraArgsString].filter(Boolean).join(' -- ')
}

const REGEXP_INDENT_LINE = /\n/g;
function indentLine(
  string,
  indentString = "  ",
  indentStringStart = indentString
) {
  return `${indentStringStart}${string.replace(
    REGEXP_INDENT_LINE,
    `\n${indentString}`
  )}`;
}

module.exports = {
  wrapJoinBashArgs,
  warpBashSubShell,
  parseCommand,
  parsePackageScript,
};
