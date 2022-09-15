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

export default function Installation({ filesNames }) {
    return <DocsPageLayout 
                title="Installation"
                data="INSTALLATION"
                filesNames={filesNames}
            />
}
