import React from 'react'
import Head from 'next/head'
import SideBar from './sideBar'
import DocBody from './docBody'

export default function DocsPageLayout({ title, data, filesNames, mdContent }) {
    return (
        <React.Fragment>
            <Head>
                <title>{title}</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            <div className="flex">
                <SideBar filesNames={filesNames} />
                <div>
                    <DocBody data={data} mdContent={mdContent} />
                </div>
            </div>
        </React.Fragment>
    )
}
