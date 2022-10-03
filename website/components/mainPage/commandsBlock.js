import { commandsList } from "../../utilities/commandsList"

function Command({ title, command }) {
    return (
        <div className="mt-4 flex">
            <p className="flex-1 typing items-center pl-0">
                <span className="text-gray-400">{`# ${title}`}</span><br/>
                <span className="text-blue-600">❯</span> {command}<br/>
            </p>
         </div>
    )
}

export default function CommandsBlock() {
    const listCommands = Object.keys(commandsList).map((title, ind) => {
        const command = commandsList[title]
        return (
            <Command key={ind}
                title={title}
                command={command}
            />
        )
    })

    return (
        <div className="flex justify-center w-full px-6 py-8 lg:h-256 lg:w-1/2">
            <div className="w-full py-12">
                <div className="coding inverse-toggle px-5 pt-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  pb-6 rounded-lg leading-normal overflow-hidden text-left">
                    <div className="top mb-2 flex">
                        <div className="h-3 w-3 bg-red-500 rounded-full"></div>
                        <div className="ml-2 h-3 w-3 bg-yellow-500 rounded-full"></div>
                        <div className="ml-2 h-3 w-3 bg-green-500 rounded-full"></div>
                    </div>

                {listCommands}
                </div>
            </div>
        </div>
    )
}
