import Head from 'next/head'
import DocsLayout from '../../components/docsLayout'

export default function Welcome() {
    return (
        <div>
            <Head>
                <title>docs</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            {/* A file path will be passed as props */}
            <DocsLayout data="A WELCOME PAGE" />
        </div>
    )
}

