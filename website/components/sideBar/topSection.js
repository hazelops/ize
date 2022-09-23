import { useState } from 'react'

import Element from './element'
import NestedSection from './nestedSection'

export default function TopSection({ title, nestedItems, id, onClickOuter, activeOuter, onClickInner, activeInner }) {
    const [isHidden, setHidden] = useState(false)

    const handleClickOuter = function() {
        setHidden(!isHidden)
        return onClickOuter(id)
    } 

    return (
        <>
            <Element 
                title={title}
                onClick={handleClickOuter}
                active={activeOuter}
                id={id}
            />
            <NestedSection
                hidden={isHidden} 
                nestedItems={nestedItems}
                onClick={onClickInner}
                active={activeInner} 
            />
        </>
    )
}
