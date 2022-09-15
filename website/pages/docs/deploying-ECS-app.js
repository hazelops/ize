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

export default function DeployingECSApp({ filesNames }) {
    return <DocsPageLayout 
                title="Deploying ECS App"
                data="Deploying ECS App"
                filesNames={filesNames}
            />
}