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

export default function Welcome({ filesNames, mdContent }) {
    return <DocsPageLayout 
                title="docs"
                data="WELCOME"
                filesNames={filesNames}
                mdContent={mdContent}
            />
}
