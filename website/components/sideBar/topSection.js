import { useState } from 'react'

import TopElement from './topElement'
import NestedSection from './nestedSection'

export default function TopSection({ title, nestedItems, currentPage }) {
    const [isHidden, setHidden] = useState(false)
    let active = false

    const handleClickOuter = function() {
        setHidden(!isHidden)
    }

    if (!currentPage[1] && title == "getting started") {
        active = true
    } else {
        const page = currentPage[currentPage.length - 1]
        const pageTitle = page.replaceAll("-", " ")
        active = nestedItems.includes(pageTitle) ? true : false
    }

    return (
        <>
            <TopElement 
                title={title}
                onClick={handleClickOuter}
                currentPage={currentPage}
                active={active}
            />
            <NestedSection
                hidden={isHidden} 
                nestedItems={nestedItems}
                currentPage={currentPage}
            />
        </>
    )
}
