import Head from 'next/head'
import DocsLayout from '../../components/docsLayout'

export default function Installation() {
    return (
        <div>
            <Head>
                <title>Deploying ECS App</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            <DocsLayout data="DEPLOYING ECS APP PAGE"/>
        </div>
    )
}