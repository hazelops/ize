import Head from 'next/head'
import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'

import { mainPageProps } from '../utilities/mainPage/mainPageProps'
import { extraData, features } from '../utilities/mainPage/featuresData'
import { commandsList } from '../utilities/mainPage/commandsList'
import { izeDescription, mainFeatures } from '../utilities/mainPage/izeData'

import IzeNavbar from '../components/izeNavbar/izeNavbar'
import TypeItAnimation from '../components/typeItAnimation'
import FeatureBlock from '../components/mainPage/featuresBlock'
import IzeFooter from '../components/izeFooter'
import CommandsBlock from '../components/mainPage/commandsBlock'
import IzeMainInfo from '../components/mainPage/izeMainInfo'

export default function Home() {
    library.add(fas)
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
                
                <section className="text-gray-600 body-font">
                    <div className="container px-5 py-24 mx-auto">
                        <div className="flex flex-wrap w-full mb-20 flex-col items-center text-center">
                            {/* <h1 id="features" className="sm:text-3xl text-2xl font-medium title-font mb-2 text-gray-900">{header}</h1> */}
                        </div>

                        <div className="flex flex-wrap -m-4">
                            <FeatureBlock title={features[0]} icon="fa-solid fa-diagram-project">
                                <p className="leading-relaxed text-base">We abstract infrastructure management and
                                    provide a clean coherent way to deploy it.</p>
                                <p className="leading-relaxed text-base">We integrate with the following tools to
                                    perform infra rollouts:</p>
                                <p className="leading-relaxed text-base">- Terraform</p>
                                <p className="leading-relaxed text-base">- Ansible*</p>
                                <p className="leading-relaxed text-base">- Cloudformation*</p>
                            </FeatureBlock>

                            <FeatureBlock title={features[1]} icon="fa-solid fa-list-check">
                                <p className="leading-relaxed text-base">We unify application deployment process and
                                    utilize naming conventions to facilitate and streamline deployments.</p>
                                <p className="leading-relaxed text-base">We allow to describe:</p>
                                <p className="leading-relaxed text-base">- ECS (currently using ecs-deploy
                                    underneath)</p>
                                <p className="leading-relaxed text-base">- k8s*</p>
                                <p className="leading-relaxed text-base">- Serverless*</p>
                            </FeatureBlock>

                            <FeatureBlock title={features[2]} icon="fa-solid fa-building-shield">
                                <p className="leading-relaxed text-base">You don’t need to setup VPN solutions to you
                                    private network, if you are just starting out.</p>
                                <p className="leading-relaxed text-base"> Also you don’t need to compromise with
                                    security.</p>
                                <p className="leading-relaxed text-base">Establish port forwarding seamlessly to any
                                    private resource via your bastion host and connect to your private resources securely.</p>
                            </FeatureBlock>

                            <FeatureBlock title={features[3]} icon="fa-solid fa-terminal">
                                <p className="leading-relaxed text-base">You can access your containers running on AWS Fargate by providing the service name.</p>
                            </FeatureBlock>

                            <FeatureBlock title={features[4]} icon="fa-solid fa-key">
                                <p className="leading-relaxed text-base">Push, Remove your secrets to/from AWS Parameter Store.</p>
                            </FeatureBlock>

                            <FeatureBlock title={features[5]} icon="fa-solid fa-seedling">
                                <p className="leading-relaxed text-base">Definitions of the environment can be stored in a toml file in the local repository.</p>
                            </FeatureBlock>
                        </div>

                        {/* <h1 className="mt-16 italic border-0 py-2 px-8">{underDev}</h1> */}
                    </div>
                </section>

                <IzeFooter />
            </div>
        </>
    )
}
