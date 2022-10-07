import Head from 'next/head'

import { mainPageProps } from '../utilities/mainPage/mainPageProps'
import { extraData, features } from '../utilities/mainPage/featuresData'
import { commandsList } from '../utilities/mainPage/commandsList'
import { izeDescription, mainFeatures } from '../utilities/mainPage/izeData'

import IzeNavbar from '../components/izeNavbar/izeNavbar'
import TypeItAnimation from '../components/typeItAnimation'
import IzeFooter from '../components/izeFooter'
import CommandsBlock from '../components/mainPage/commandsBlock'
import IzeMainInfo from '../components/mainPage/izeMainInfo'
import FeaturesBlock from '../components/mainPage/featuresBlock'

export default function Home() {
    const { pageTitle, description, previewImage } = mainPageProps
    return (
        <>
            <Head>
                <title>{pageTitle}</title>
                <link rel="icon" href="/favicon.ico" />
                <meta name="viewport" content="width=device-width, initial-scale=1" />
                <meta charSet="utf-8" />

                <meta name="description" content={description} />
                <meta property="og:title" content={pageTitle} key="ogtitle" />
                <meta property="og:description" content={description} key="ogdesc" />

                {/* Open Graph */}
                {/*<meta property="og:url" content={props.currentURL} key="ogurl" />*/}
                <meta property="og:image" content={previewImage} key="ogimage" />
                <meta property="og:site_name" content={pageTitle} key="ogsitename" />

            </Head>

            <div className="flex flex-col">
                <section className="bg-white dark:bg-gray-800">
                    <IzeNavbar />

                    <div className="lg:flex">
                        <div className="flex justify-center w-full px-6 py-6 lg:h-128 lg:w-1/2">
                            <div className="max-w-xl">
                                <TypeItAnimation />

                                <IzeMainInfo izeDescription={izeDescription} mainFeatures={mainFeatures} />
                            </div>
                        </div>

                        <CommandsBlock commandsList={commandsList} />
                    </div>
                </section>
                
                <FeaturesBlock extraData={extraData} features={features} />

                <IzeFooter />
            </div>
        </>
    )
}
