import styles from './sideBar.module.css'

export default function Element({ title,  onClick, currentPage }) {
    let color

    if (!currentPage[1] && title == "getting started") {
        color = "selectedElement"
    } else {
        const fileName = currentPage[currentPage.length - 1].replaceAll("-", " ")
        color = fileName == title ? "selectedElement" : ""
    }
   
    return (
            <div
                className={`${styles.topElement} w-fit py-2 mt-2 cursor-pointer ${color}`}
                onClick={onClick}
            >
                {title}
            </div>
    )
}
