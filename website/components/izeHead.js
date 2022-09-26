import Head from 'next/head'

export default function IzeHead(props) {
    return (
            <Head>
                <title>{props.title}</title>
                <link rel="icon" href="/favicon.ico"/>
                <meta name="viewport" content="width=device-width, initial-scale=1"/>
                <meta charSet="utf-8"/>
            </Head>
    )
}
