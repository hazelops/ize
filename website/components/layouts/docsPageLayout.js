import { useRouter } from 'next/router'

import IzeHead from '../izeHead'
import SideBar from '../sideBar/sideBar'
import DocsNavbar from '../docsNavbar/docsNavbar'

export default function DocsPageLayout({ children, title, filesNames }) {
    const router = useRouter()
    const currentPage = router.pathname.split("/")
    currentPage.shift()

    return (
        <>
            <IzeHead title={title} />

            <div className="flex w-full">
                <SideBar filesNames={filesNames} currentPage={currentPage} />
                <div className="flex flex-col">
                   <DocsNavbar />
                    <div className="flex w-full justify-center pl-10 pt-10">
                        {children}
                    </div> 
                </div>
            </div>
        </>
    )
}
