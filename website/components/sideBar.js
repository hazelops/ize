import Ize from "./ize"
import React, { useState } from 'react'
import Link from 'next/link'
import { sideBarMenu } from '../utilities/sideBarMenu'

function TopElement(props) {
    return (
            <div id={props.title} className="flex items-center px-4 py-2 mt-5 text-gray-600 rounded-md hover:bg-gray-200 transition-colors duration-300 transform cursor-pointer" onClick={props.onClick}>
                <span className="mx-4 font-medium capitalize">{props.title}</span>
            </div>
    )
}

function NestedMenu(props) {
    if (props.hidden) {
        return null
    }

    const nestedList = props.nestedItems.map(el => {
        const pathName = el.slice().replaceAll(" ", "-")
        let route = pathName == "welcome"? "" : pathName
        return <Link 
                key={props.nestedItems.indexOf(el)} 
                href={`/docs/${route}`}>
                    <a><TopElement title={el}/></a>
                </Link>
    })

    return (
        <div className="flex flex-col justify-between flex-1 ml-10">
            {nestedList}
        </div>
    )
}

function MenuElement(props) {
    const [isHidden, setHidden] = useState(false)

    const handleClick = function() {
        setHidden(!isHidden)
    } 

    return (
        <React.Fragment>
            <TopElement title={props.title} onClick={handleClick} />
            <NestedMenu hidden={isHidden} nestedItems={props.nestedItems} />
        </React.Fragment>
        
    )
}
//------------------------------------------------------------------------------------------------------------------

export default function SideBar(props) {
    const { mainMenu, seeAlso } = sideBarMenu

    const menuList = mainMenu.map(el => {
        return (
            <MenuElement
                key={mainMenu.indexOf(el)}
                title={el.title}
                nestedItems={el.nestedItems}
             />
        )
    })

    return (
        <div className="flex flex-col w-64 h-screen px-4 py-8 bg-white border-r">
            <div>
               <Ize /> 
            </div>

            <div className="flex flex-col justify-between flex-1">
                <nav>
                    {menuList}
                    <hr className="my-6 border-gray-200" />
                    <TopElement
                        title={seeAlso.title}
                    />
                </nav>
            </div>
        </div>
    )
}