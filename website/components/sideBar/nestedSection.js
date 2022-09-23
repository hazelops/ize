import Link from 'next/link'

import styles from './sideBar.module.css'

function NestedElement({ title, active }) {
    const color = active ? "selectedElement" : ""

    return (
        <div
            className={`${styles.topElement} w-fit py-2 mt-2 cursor-pointer ${color}`}
        >
            {title}
        </div>
    )
}

//---------------------------------------------------------- 

export default function NestedSection({ hidden, nestedItems, currentPage }) {
    if (hidden) {
        return null
    }

    const nestedList = nestedItems.map(el => {
        const ind = nestedItems.indexOf(el)
        let active = false

        const pathName = el.slice().replaceAll(" ", "-")
        let route = pathName == "welcome" ? "" : pathName

        if (!currentPage[1] && el == "welcome") {
            active = true
        } else {
            const page = currentPage[currentPage.length - 1]
            active = page == pathName ? true : false
        }

        return <Link key={ind} href={`/docs/${route}`}>
                    <a>
                        <NestedElement 
                            title={el}
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
