import { titles } from '../../utilities/docsNavbarMenu'
import GitHubIcon from '../gitHubIcon'

import styles from './docsNavbar.module.css'

function NavButton({ title }) {
    return (
        <div className={`${styles.button} px-2 py-1.5 capitalize`} role="button">
            {title}
        </div>
    )
}

export default function DocsNavbar() {
    const listButtons = titles.map(title => {
        const ind = titles.indexOf(title)
        return <NavButton key={ind} title={title} />
    })

    return (
        <nav className="flex justify-between">
            <div className={`${styles.outer} flex justify-between items-center w-1/3 px-10 pt-5`}>
                {listButtons}
            </div>

            <div className={`${styles.iconOuter} pt-5 self-center`}>
                <GitHubIcon />
            </div>
        </nav>
    )
}
