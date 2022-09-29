import Ize from '../ize'
import { headers } from '../../utilities/welcomePageHeaders'

import styles from './welcomeContent.module.css'

export default function WelcomeContent() {
    const listHeaders = Object.keys(headers).map((header, ind) => {
        return <h2 key={ind} className="">{header}</h2>
    })

    return (
        <>
            <header className={styles.header}>
                <h1 className="text-3xl">
                    Welcome to 
                    <span className="px-4">
                        <Ize /> 
                    </span>
                    docs!
                </h1>
            </header> 

            <section>
                
            </section>
        </>
    )
}
