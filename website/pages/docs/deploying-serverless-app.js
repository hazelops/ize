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

export default function DeployingServerlessApp({ filesNames, mdContent }) {
    return <DocsPageLayout 
                title="Deploying Serverless App"
                data="Deploying Serverless App"
                filesNames={filesNames}
                mdContent={mdContent}
            />
}
