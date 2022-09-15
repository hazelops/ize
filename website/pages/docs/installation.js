import Head from 'next/head'
import DocsLayout from '../../components/docsLayout'

export default function Installation() {
    return (
        <div>
            <Head>
                <title>Installation</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>
            <DocsLayout data="AN INSTALLATION PAGE"/>
        </div>
    )
}
