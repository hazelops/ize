import { useState } from 'react'
import { installationMenu } from '../../utilities/installationMenu'

function TabButton({ field, id, onClick, active }) {
    const border = active == id ? "border-b-2 border-blue-500" : ""

    return (
        <button 
            className={`h-10 px-4 py-2 -mb-px text-sm text-center text-blue-600 bg-transparent sm:text-base whitespace-nowrap focus:outline-none hover:bg-gray-200 ${border}`}
            onClick={() => onClick(id)}
        >
            {field}
        </button>
    )
}

function TabContent({ children }) {
    return (
        <div>{children}</div>
    )
}

//---------------------------------------------------------

export default function Tab() {
    const [active, setActive] = useState(0)

    const handleClick = function(id) {
        return setActive(id)
    }

    const listButtons = installationMenu.map(el => {
        const index = installationMenu.indexOf(el)
        return (
            <TabButton key={index}
                field={el} 
                onClick={handleClick} 
                id={index}
                active={active} 
            />
        )
    })

    return (
            <div className="flex flex-col w-2/3 items-center">
                <div className="inline-flex w-1/2 justify-between border-b border-gray-200 h-fit">
                    {listButtons}
                </div>
                <TabContent></TabContent>
            </div> 
    )
}
