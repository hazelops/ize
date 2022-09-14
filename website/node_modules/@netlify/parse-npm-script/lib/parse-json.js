const fs = require('fs')
const { promisify } = require('util')
const readFileAsync = promisify(fs.readFile)

async function parseJson(packagePath, script) {
  const pkgString = await readFileAsync(packagePath, 'utf-8')
  const pkg = JSON.parse(pkgString)
  return pkg
}

module.exports = parseJson
