import { useState } from 'react'

function TabButton(props) {
    const [selected, setSelected] = useState()
    const border = selected ? "border-b-2 border-blue-500" : ""

    return (
        <button 
            className={`h-10 px-4 py-2 -mb-px text-sm text-center text-blue-600 bg-transparent sm:text-base whitespace-nowrap focus:outline-none ${border}`}
            onClick={props.onClick}
        >
            {props.platform}
        </button>
    )
}

export default function Tab() {
    const handleClick = function() {
        return
    }

    return (
            <div className="inline-flex w-2/3 justify-center">
                <div className="inline-flex w-1/2 justify-between border-b border-gray-200 h-fit">
                    <TabButton platform="MacOS" />
                    <TabButton platform="Ubuntu" />
                    <TabButton platform="Source" />
                    <TabButton platform="Autocomplete" />
                </div>
            </div> 
    )
}
