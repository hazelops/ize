import React, { useState } from 'react'
import Ize from "../ize"
import Link from 'next/link'

import { sideBarMenu } from '../../utilities/sideBarMenu'
import styles from './sideBar.module.css'

function TopElement({ title, onClick }) {
    return (
            <div
                className={`${styles.topElement} w-fit py-2 mt-3 cursor-pointer`}
                onClick={onClick}
                >
                <span className="font-medium capitalize">{title}</span>
            </div>
    )
}

function NestedElements({ hidden, nestedItems }) {
    if (hidden) {
        return null
    }

    const nestedList = nestedItems.map(el => {
        const pathName = el.slice().replaceAll(" ", "-")
        let route = pathName == "welcome"? "" : pathName
        return <Link 
                key={nestedItems.indexOf(el)} 
                href={`/docs/${route}`}
                >
                    <a><TopElement title={el}/></a>
                </Link>
    })

    return (
        <div className="flex flex-col justify-between flex-1 ml-7">
            {nestedList}
        </div>
    )
}

function MenuElement({ title, nestedItems }) {
    const [isHidden, setHidden] = useState(false)

    const handleClick = function() {
        setHidden(!isHidden)
    } 

    return (
        <React.Fragment>
            <TopElement title={title} onClick={handleClick} />
            <NestedElements hidden={isHidden} nestedItems={nestedItems} />
        </React.Fragment>
        
    )
}

//------------------------------------------------------------------------------------------------------------------

export default function SideBar({ filesNames }) {
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
        <div className={`${styles.outer} flex flex-col w-1/6 px-5 py-3`}>
                <nav className="">
                    <div className="flex mb-5">
                        <Ize /> 
                        <div className="text-2xl font-bold text-blue-600 lg:text-3xl ml-5">docs</div>
                    </div>
                    {menuList}
                    <hr className="my-6 border-gray-200" />
                    <MenuElement
                        title={seeAlso.title}
                        nestedItems={filesNames}
                    />
                </nav>
        </div>
    )
}
