import React from 'react'
import IzeHead from './izeHead'
import SideBar from './sideBar'

export default function DocsPageLayout({ children, title, filesNames }) {
    return (
        <React.Fragment>
            <IzeHead title={title} />

            <div className="flex">
                <SideBar filesNames={filesNames} />
                <div>
                    {children}
                </div>
            </div>
        </React.Fragment>
    )
}
