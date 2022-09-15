import Head from 'next/head'
import DocsLayout from '../../components/docsLayout'
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
    return (
        <div>
            <Head>
                <title>docs</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>
            <DocsLayout 
                data="A WELCOME PAGE"
                filesNames={filesNames} 
            />
        </div>
    )
}
