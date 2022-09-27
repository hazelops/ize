import DocsPageLayout from '../../components/layouts/docsPageLayout'
import Tab from '../../components/tab/tab'
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
                <Tab />
            </DocsPageLayout>
}
