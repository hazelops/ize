import Head from 'next/head'
import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'

import { mainPageProps } from '../utilities/mainPageProps'
import { features } from '../utilities/features'

import IzeNavbar from '../components/izeNavbar/izeNavbar'
import TypeItAnimation from '../components/typeItAnimation'
import FeatureBlock from '../components/mainPage/featureBlock'
import IzeFooter from '../components/izeFooter'
import CommandsBlock from '../components/mainPage/commandsBlock'

export async function getStaticProps() {
    const { pageTitle, description, previewImage } = mainPageProps
    const titles = features
    return {
        props: {
            pageTitle,
            description,
            previewImage,
            titles
        }
    }
}

export default function Home({ pageTitle, description, previewImage, titles }) {
    library.add(fas)
    return (
        <>
            <Head>
                <title>{pageTitle}</title>
                <link rel="icon" href="/favicon.ico"/>
                <meta name="viewport" content="width=device-width, initial-scale=1"/>
                <meta charSet="utf-8"/>

                <meta name="description" content={description}/>
                <meta property="og:title" content={pageTitle} key="ogtitle"/>
                <meta property="og:description" content={description} key="ogdesc"/>

                {/* Open Graph */}
                {/*<meta property="og:url" content={props.currentURL} key="ogurl" />*/}
                <meta property="og:image" content={previewImage} key="ogimage"/>
                <meta property="og:site_name" content={pageTitle} key="ogsitename"/>

            </Head>

            <div className="flex flex-col">
                <section className="bg-white dark:bg-gray-800">
                    <IzeNavbar />

                    <div className="lg:flex">
                        <div className="flex justify-center w-full px-6 py-6 lg:h-128 lg:w-1/2">
                            <div className="max-w-xl">
                                <TypeItAnimation />

                                <p className="mt-5 text-sm text-gray-500 dark:text-gray-400 lg:text-base">
                                    An opinionated deployment tool for infrastructure and code. The main goal is to co-join
                                    operational
                                    tasks into one easy-to-use tool. It is written in Go provides a robust abstraction level
                                    on top of
                                    common orchestration and code deployment systems.</p>
                                <div className="grid gap-6 mt-8 sm:grid-cols-2">
                                    <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                        <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                                d="M5 13l4 4L19 7"/>
                                        </svg>

                                        <span>AWS</span>
                                    </div>

                                    <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                        <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                                d="M5 13l4 4L19 7"/>
                                        </svg>

                                        <span>Terraform</span>
                                    </div>

                                    <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                        <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                                d="M5 13l4 4L19 7"/>
                                        </svg>

                                        <span>Docker</span>
                                    </div>

                                    <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                        <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                                d="M5 13l4 4L19 7"/>
                                        </svg>

                                        <span>Serverless</span>
                                    </div>

                                    <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                        <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                                d="M5 13l4 4L19 7"/>
                                        </svg>

                                        <span>ECS</span>
                                    </div>

                                    <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                        <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                                d="M5 13l4 4L19 7"/>
                                        </svg>

                                        <span>CI/CD</span>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <CommandsBlock />
                    </div>


                </section>
                <section className="text-gray-600 body-font">
                    <div className="container px-5 py-24 mx-auto">
                        <div className="flex flex-wrap w-full mb-20 flex-col items-center text-center">
                            <h1 id="features" className="sm:text-3xl text-2xl font-medium title-font mb-2 text-gray-900">Features</h1>
                        </div>
                        <div className="flex flex-wrap -m-4">
                            <FeatureBlock title={titles[0]} icon="fa-solid fa-diagram-project">
                                    <p className="leading-relaxed text-base">We abstract infrastructure management and
                                        provide a clean coherent way to deploy it.</p>
                                    <p className="leading-relaxed text-base">We integrate with the following tools to
                                        perform infra rollouts:</p>
                                    <p className="leading-relaxed text-base">- Terraform</p>
                                    <p className="leading-relaxed text-base">- Ansible*</p>
                                    <p className="leading-relaxed text-base">- Cloudformation*</p>
                            </FeatureBlock>

                            <FeatureBlock title={titles[1]} icon="fa-solid fa-list-check">
                                    <p className="leading-relaxed text-base">We unify application deployment process and
                                        utilize naming conventions to facilitate and streamline deployments.</p>
                                    <p className="leading-relaxed text-base">We allow to describe:</p>
                                    <p className="leading-relaxed text-base">- ECS (currently using ecs-deploy
                                        underneath)</p>
                                    <p className="leading-relaxed text-base">- k8s*</p>
                                    <p className="leading-relaxed text-base">- Serverless*</p>
                            </FeatureBlock>
                        
                            <FeatureBlock title={titles[2]} icon="fa-solid fa-building-shield">
                                    <p className="leading-relaxed text-base">You don’t need to setup VPN solutions to you
                                        private network, if you are just starting out.</p>
                                    <p className="leading-relaxed text-base"> Also you don’t need to compromise with
                                        security.</p>
                                    <p className="leading-relaxed text-base">Establish port forwarding seamlessly to any
                                        private resource via your bastion host and connect to your private resources securely.</p>   
                            </FeatureBlock>            
                            
                            <FeatureBlock title={titles[3]} icon="fa-solid fa-terminal">
                                    <p className="leading-relaxed text-base">You can access your containers running on AWS Fargate by providing the service name.</p>        
                            </FeatureBlock>

                            <FeatureBlock title={titles[4]} icon="fa-solid fa-key">
                                    <p className="leading-relaxed text-base">Push, Remove your secrets to/from AWS Parameter Store.</p>   
                            </FeatureBlock>            
                        
                            <FeatureBlock title={titles[5]} icon="fa-solid fa-seedling">
                                    <p className="leading-relaxed text-base">Definitions of the environment can be stored in a toml file in the local repository.</p>    
                            </FeatureBlock>            
                        </div>
                        <h1 className="mt-16 italic border-0 py-2 px-8">*Currently under development</h1>
                    </div>
                </section>

                <IzeFooter />
            </div>
        </>
    )
}
