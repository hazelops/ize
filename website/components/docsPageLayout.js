import React from 'react'
import IzeHead from './izeHead'
import SideBar from './sideBar/sideBar'

export default function DocsPageLayout({ children, title, filesNames }) {
    return (
        <React.Fragment>
            <IzeHead title={title} />

            <div className="flex">
                <SideBar filesNames={filesNames} />
                <div className="flex w-full justify-center">
                    {children}
                </div>
            </div>
        </React.Fragment>
    )
}
