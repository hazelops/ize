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

export default function DeployingECSApp({ filesNames, mdContent }) {
    return <DocsPageLayout 
                title="Deploying ECS App"
                filesNames={filesNames}
            >
                <DocBody mdContent={mdContent} />
            </DocsPageLayout>
}