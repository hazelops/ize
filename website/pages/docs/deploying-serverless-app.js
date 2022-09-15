import DocsPageLayout from '../../components/docsPageLayout'
import { readFilesNames } from '../../utilities/readFilesNames'

export async function getServerSideProps() {
    const filesNames = await readFilesNames()
    return {
        props: {
         filesNames
        }
    }
}

export default function DeployingServerlessApp({ filesNames }) {
    return <DocsPageLayout 
                title="Deploying Serverless App"
                data="Deploying Serverless App"
                filesNames={filesNames}
            />
}
