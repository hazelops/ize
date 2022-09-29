import Ize from '../ize'

import styles from './welcomeContent.module.css'

export default function WelcomeContent() {
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
