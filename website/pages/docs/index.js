import DocsPageLayout from '../../components/layouts/docsPageLayout'
import DocBody from '../../components/docBody'
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
                filesNames={filesNames}
            >
                <DocBody mdContent={mdContent} />
            </DocsPageLayout>
}
