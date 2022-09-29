import { sideBarMenu } from '../../utilities/sideBarMenu'
import Ize from "../ize"
import TopSection from './topSection'

import styles from './sideBar.module.css'

export default function SideBar({ filesNames, currentPage }) {
    const { mainMenu, seeAlso } = sideBarMenu

    const menuList = mainMenu.map(el => {
        const ind = mainMenu.indexOf(el)
        return (
            <TopSection key={ind}
                title={el.title}
                nestedItems={el.nestedItems}
                currentPage={currentPage}
             />
        )
    })

    return (
        <nav className={`${styles.outer} flex flex-col px-5 pt-7`}>
            <header className="flex mb-5">
                <Ize /> 
                <div className="text-3xl font-bold text-blue-600 ml-5">docs</div>
            </header>

            {menuList}

            <hr className="my-6 border-gray-200" />

            <TopSection
                title={seeAlso.title}
                nestedItems={filesNames}
                currentPage={currentPage}
            />
        </nav>
    )
}
