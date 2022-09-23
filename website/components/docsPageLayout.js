import { useRouter } from 'next/router'

import IzeHead from './izeHead'
import SideBar from './sideBar/sideBar'

export default function DocsPageLayout({ children, title, filesNames }) {
    const router = useRouter()
    const currentPage = router.pathname.split("/")
    currentPage.shift()

    return (
        <>
            <IzeHead title={title} />

            <div className="flex">
                <SideBar filesNames={filesNames} currentPage={currentPage} />
                <div className="flex w-full justify-center">
                    {children}
                </div>
            </div>
        </>
    )
}
