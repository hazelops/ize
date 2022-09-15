import DocsPageLayout from '../../components/docsPageLayout'
import { readFilesNames, fetchContent } from '../../utilities/docsGlobalProps'

export async function getStaticProps() {
    const filesNames = await readFilesNames()
    const mdContent = fetchContent()
    return {
        props: {
         filesNames,
         mdContent
        }
    }
}

export default function Installation({ filesNames, mdContent }) {
    return <DocsPageLayout 
                title="Installation"
                data="INSTALLATION"
                filesNames={filesNames}
                mdContent={mdContent}
            />
}
