import styles from './sideBar.module.css'

export default function TopElement({ title,  onClick, active }) {
    let color = active ? "selectedElement" : ""

    return (
        <div
            className={`${styles.topElement} w-fit py-2 mt-2 cursor-pointer ${color}`}
            onClick={onClick}
        >
            {title}
        </div>
    )
}
