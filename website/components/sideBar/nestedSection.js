import Link from 'next/link'

import Element from './element'

export default function NestedSection({ hidden, nestedItems, currentPage }) {
    if (hidden) {
        return null
    }

    const nestedList = nestedItems.map(el => {
        const ind = nestedItems.indexOf(el)

        const pathName = el.slice().replaceAll(" ", "-")
        let route = pathName == "welcome" ? "" : pathName

        return <Link key={ind} href={`/docs/${route}`}>
                    <a>
                        <Element 
                            title={el}
                            currentPage={currentPage}
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
