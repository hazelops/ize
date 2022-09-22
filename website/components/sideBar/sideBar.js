import React, { useState } from 'react'
import Ize from "../ize"
import Link from 'next/link'

import { sideBarMenu } from '../../utilities/sideBarMenu'
import styles from './sideBar.module.css'

function TopElement({ title, onClick, active, id }) {
    const color = active == id ? "selectedElement" : ""
    return (
            <div
                className={`${styles.topElement} w-fit py-2 mt-2 cursor-pointer ${color}`}
                onClick={onClick}
            >
                {title}
            </div>
    )
}

function NestedElements({ hidden, nestedItems, active, onClick }) {
    if (hidden) {
        return null
    }

    const nestedList = nestedItems.map(el => {
        const ind = nestedItems.indexOf(el)

        const handleClick = function() {
            return onClick(ind)
        }

        const pathName = el.slice().replaceAll(" ", "-")
        let route = pathName == "welcome"? "" : pathName

        return <Link key={ind} href={`/docs/${route}`}>
                    <a>
                        <TopElement 
                            title={el}
                            id={ind} 
                            onClick={handleClick}
                            active={active}
                        />
                    </a>
                </Link>
    })

    return (
        <div className="flex flex-col justify-between flex-1 ml-7">
            {nestedList}
        </div>
    )
}

function MenuElement({ title, nestedItems, id, onClickOuter, activeOuter, onClickInner, activeInner }) {
    const [isHidden, setHidden] = useState(false)

    const handleClickOuter = function() {
        setHidden(!isHidden)
        return onClickOuter(id)
    } 

    return (
        <React.Fragment>
            <TopElement 
                title={title}
                onClick={handleClickOuter}
                active={activeOuter}
                id={id}
            />
            <NestedElements 
                hidden={isHidden} 
                nestedItems={nestedItems}
                onClick={onClickInner}
                active={activeInner} 
            />
        </React.Fragment>
        
    )
}

//------------------------------------------------------------------------------------------------------------------

export default function SideBar({ filesNames }) {
    const { mainMenu, seeAlso } = sideBarMenu

    const [activeOuter, setActiveOuter] = useState(0)
    const [activeInner, setActiveInner] = useState(null)
    // get unique ids (array?)

    const handleClickOuter = function(id) {
        setActiveOuter(id)
    }

    const handleClickInner = function(id) {
        setActiveInner(id)
    }

    const menuList = mainMenu.map(el => {
        const ind = mainMenu.indexOf(el)
        return (
            <MenuElement key={ind}
                id={ind}
                title={el.title}
                nestedItems={el.nestedItems}
                onClickOuter={handleClickOuter}
                activeOuter={activeOuter}
                onClickInner={handleClickInner}
                activeInner={activeInner}
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
                        id={2}
                        onClickOuter={handleClickOuter}
                        activeOuter={activeOuter}
                        onClickInner={handleClickInner}
                        activeInner={activeInner}
                    />
                </nav>
        </div>
    )
}
