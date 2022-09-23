import { useState } from 'react'

import Element from './element'
import NestedSection from './nestedSection'

export default function TopSection({ title, nestedItems, currentPage }) {
    const [isHidden, setHidden] = useState(false)
    const [isActive, setActive] = useState(false)

    const handleClickOuter = function() {
        setHidden(!isHidden)
    }

    return (
        <>
            <Element 
                title={title}
                onClick={handleClickOuter}
                currentPage={currentPage}
            />
            <NestedSection
                hidden={isHidden} 
                nestedItems={nestedItems}
                currentPage={currentPage}
            />
        </>
    )
}
