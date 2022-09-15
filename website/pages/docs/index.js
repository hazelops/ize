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

export default function Welcome({ filesNames }) {
    return <DocsPageLayout 
                title="docs"
                data="WELCOME"
                filesNames={filesNames}
            />
}
