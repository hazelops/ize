import DocsPageLayout from '../../components/layouts/docsPageLayout'
import WelcomeContent from '../../components/welcomePageContent/welcomeContent'
import { readFilesNames } from '../../utilities/docsGlobalProps'

export async function getStaticProps() {
    const filesNames = await readFilesNames()
    return {
        props: {
         filesNames
        }
    }
}

export default function Welcome({ filesNames, mdContent }) {
    return <DocsPageLayout 
                title="docs"
                filesNames={filesNames}
            >
               <WelcomeContent />
            </DocsPageLayout>
}
