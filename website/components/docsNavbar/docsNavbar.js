import { titles } from '../../utilities/docsNavbarMenu'
import styles from './docsNavbar.module.css'

function NavButton({ title }) {
    return (
        <button className="capitalize">
            {title}
        </button>
    )
}

export default function DocsNavbar() {
    const listButtons = titles.map(title => {
        return <NavButton title={title} />
    })

    return (
        <nav className={`${styles.outer} flex justify-between w-1/3 px-10`}>
            {listButtons}
        </nav>
    )
}
