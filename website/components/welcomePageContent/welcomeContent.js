import React from 'react'
import Ize from '../ize'

import styles from './welcomeContent.module.css'

function Quickstart({ data }) {
    const { title, content } = data
    return (
        <>
            <h2 className={`${styles.contentHeader} pt-8`}>{title}</h2>
            <span>BLOCK BLOCK BLOCK</span>
        </>
    )
}

function WhatIsIze({ data }) {
    const { title, content } = data
    const listSubHeaders = content.map((subHeader, ind) => {
        return (
            <React.Fragment key={ind}>
                <h3 className={styles.contentSubHeader}>{subHeader}</h3>
                <div className={styles.content}>
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque at urna ultricies, iaculis mi accumsan, faucibus nisl. Ut id ullamcorper nunc. Duis dignissim tempor tortor, id blandit dui volutpat sit amet. Cras ornare lectus vel mi aliquet tristique. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut in massa metus. Nulla at quam sem. Donec a tincidunt ipsum, vitae laoreet purus. Vestibulum commodo, enim quis imperdiet consectetur, risus ligula cursus mi, eget elementum neque lectus vel diam. Integer consectetur euismod justo eleifend eleifend. Cras maximus interdum cursus. Etiam consectetur leo sit amet enim vulputate elementum.
                </div>
            </React.Fragment>
        )
    })

    return (
        <>
            <h2 className={`${styles.contentHeader} pt-10`}>{title}</h2>
            {listSubHeaders}
        </>
    )
}

// --------------------------------------------

export default function WelcomeContent({ headers }) {
    const [ quickstart, whatIsIze ] = headers
    return (
        <section className={styles.outer}>
            <header className={styles.header}>
                <h1>
                    Welcome to
                    <span className="px-4">
                        <Ize />
                    </span>
                    docs!
                </h1>
            </header>

            <Quickstart data={quickstart} />
            <WhatIsIze data={whatIsIze} />
        </section>
    )
}
