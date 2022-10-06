import React from 'react'
import Link from 'next/link'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

import Ize from '../ize'
import Chevron from '../chevron'
import styles from './welcomeContent.module.css'

function QuickstartBlock({ blockHeader, icon, text }) {
    return (
        <div className="flex flex-col items-center border border-gray-200 rounded-lg p-4">
            <div className="flex items-center mb-4">
                <div className="inline-block w-[1rem] h-[1rem] text-blue-600 mx-2">
                    <FontAwesomeIcon icon={icon} />
                </div> 
                <div className={styles.blockHeader}>{blockHeader}</div>
                <div className="inline-block w-[1rem] h-[1rem] text-blue-600 mx-2">
                    <FontAwesomeIcon icon={icon} />
                </div> 
            </div>

            <div className="mb-4">
                {text}
            </div>

            <div className={styles.link}>
                <Link href="#">
                    <a>
                        Continue
                        <span className="ml-2">
                            <Chevron />
                        </span>
                    </a>
                </Link>
            </div> 
        </div>
    )
}

function Quickstart({ data }) {
    const { title, content } = data
    const blockHeaders = Object.keys(content)
    const listBlocks = blockHeaders.map((blockHeader, ind) => {
        const icon = content[blockHeader].icon
        const text = content[blockHeader].text.concat(".")
        return (
            <QuickstartBlock key={ind} 
                blockHeader={blockHeader}
                icon={icon}
                text={text}
            />
        )
    })

    return (
        <div className="w-1/2">
            <h2 className={`${styles.contentHeader} mt-4`}>{title}</h2>
            <div className="flex justify-between mt-9">
                {listBlocks}
            </div>
        </div>
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
        <div>
            <h2 className={`${styles.contentHeader} mt-8`}>{title}</h2>
            {listSubHeaders}
        </div>
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
