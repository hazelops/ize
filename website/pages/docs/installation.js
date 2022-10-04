import DocsPageLayout from '../../components/layouts/docsPageLayout'
import InstallationTab from '../../components/tab/InstallationTab'
import { readFilesNames } from '../../utilities/docsGlobalProps'

export async function getStaticProps() {
    const filesNames = await readFilesNames()
    return {
        props: {
         filesNames
        }
    }
}

export default function Installation({ filesNames }) {
    return <DocsPageLayout 
                title="Installation"
                filesNames={filesNames}
            >
                <InstallationTab />
            </DocsPageLayout>
}
