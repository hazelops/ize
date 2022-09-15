export const readFilesNames = function() {
    const fs = require('fs')
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
