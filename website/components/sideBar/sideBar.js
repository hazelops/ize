import { useState } from 'react'

import { sideBarMenu } from '../../utilities/sideBarMenu'
import Ize from "../ize"
import TopSection from './topSection'

import styles from './sideBar.module.css'

export default function SideBar({ filesNames }) {
    const { mainMenu, seeAlso } = sideBarMenu

    const [activeOuter, setActiveOuter] = useState(0)
    const [activeInner, setActiveInner] = useState(null)
    const [innerIds, setInnerIds] = useState([])

    const handleClickOuter = function(id) {
        setActiveOuter(id)
    }

    const handleClickInner = function(id) {
        setActiveInner(id)
    }

    const menuList = mainMenu.map(el => {
        const ind = mainMenu.indexOf(el)
        return (
            <TopSection key={ind}
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
                    <TopSection
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
