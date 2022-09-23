import styles from './sideBar.module.css'

export default function Element({ title, onClick, active, id }) {
    const color = active == id ? "selectedElement" : ""
    return (
            <div
                className={`${styles.topElement} w-fit py-2 mt-2 cursor-pointer ${color}`}
                onClick={onClick}
            >
                {title}
            </div>
    )
}
