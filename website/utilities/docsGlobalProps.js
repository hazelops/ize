import fs from 'fs'
import path from 'path'
import md from 'markdown-it'

export const readFilesNames = function() {
    return new Promise(resolve => {
        fs.readdir("./docs", (err, files) => {
            if (err) {
                console.log("error")
            }
            const filesNames = files.map(file => file.replace(".md", ""))
            resolve(filesNames)
        })
    })
} 

export const fetchContent = function() {
    const filePath = path.join(process.cwd(), 'docs', 'doc1.md')
    const fileContents = fs.readFileSync(filePath, 'utf8')
    const result = md().render(fileContents)
    return result
}
