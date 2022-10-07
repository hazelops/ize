import DocsPageLayout from '../../components/layouts/docsPageLayout'
import WelcomeContent from '../../components/welcomePageContent/welcomeContent'

import { readFilesNames } from '../../utilities/docsGlobalProps'
import { headers } from '../../utilities/welcomePageHeaders'


export async function getStaticProps() {
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
                filesNames={filesNames}
            >
               <WelcomeContent headers={headers} />
            </DocsPageLayout>
}
