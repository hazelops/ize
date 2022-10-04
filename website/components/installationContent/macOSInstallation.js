import React from "react";

import CodeTemplate from "../codeTemplate";
import Chevron from "../chevron"
import { data } from "../../utilities/installationPage/macOS";

// function Print({ content }) {
//     const descriptions = Object.keys(content)
//     descriptions.map((el, ind) => {
//         if (!content[el]) {
//             return (
//                 <React.Fragment key={ind}>
//                     <li className="leading-relaxed text-base">{el}</li>
//                 </React.Fragment>
//             )
//         }
//         const listCommands = function() {
//             const commands = content[el]
//             commands.map((command, ind) => {
//                 return (
//                     <>
//                       <Chevron /> {command}  
//                     </>
                    
//                 )
//             })
//         }
//         return (
//             <React.Fragment key={ind}>
//                 <li className="leading-relaxed text-base">{el}</li>
//                 <CodeTemplate>
//                     {listCommands}
//                 </CodeTemplate>
//             </React.Fragment>
//         )
//     })
// }

function ListCommands({ commands }) {
    return commands.map((el, ind) => {
        return (
            <div key={ind}>
                <Chevron /> {el}
            </div>
        )
    })
}

export default function MacOSInstallation() {
    const { install, update } = data
    const installCommands = data.install.content["Run the following commands:"]
    return (
        <>
            <h3 className="text-xl font-medium mb-2 text-gray-700">
                {install.header}
            </h3>
            <h5 className="leading-relaxed text-base">
                1.1 Install
                <a target="_blank" href="https://brew.sh/">Homebrew</a>
            </h5>                    
            <h5 className="leading-relaxed text-base">
                1.2 Run the following commands:
            </h5>
            <CodeTemplate>
                <ListCommands commands={installCommands} />
            </CodeTemplate>
            <p className="leading-relaxed text-base">
                Now you can run ize from command shell by typing <code>ize</code> in console.
            </p>


            <h5 className="text-xl font-medium mb-2 text-gray-700">
                {update.header}
            </h5>
            <p className="leading-relaxed text-base">
                To update ize you should:
            </p>
            <h5 className="leading-relaxed text-base">
                2.1 Uninstall previous version:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span> 
                brew uninstall ize <br/>
            </CodeTemplate>
            <h5 className="leading-relaxed text-base">
                2.2 Update version in brew repo:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span> 
                brew tap hazelops/ize <br/>
            </CodeTemplate>
            <h5 className="leading-relaxed text-base">
                2.3 Install:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span> 
                brew install ize <br/>
            </CodeTemplate>
        </>
    )
}
