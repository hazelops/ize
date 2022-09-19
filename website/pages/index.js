// https://tailblocks.cc/
// https://merakiui.com/
// https://www.tailwind-kit.com/components#elements
import Head from 'next/head'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'

import IzeNavbar from '../components/izeNavbar'
import TypeItAnimation from '../components/typeItAnimation'

library.add(fas)

export default function Home() {
    let props = {
        pageTitle: "ize: Opinionated Infra Tool",
        description: "Opinionated Infra Tool",
        previewImage: "/social-preview.png"
    }

    return (
        <div className="flex flex-col">
            <Head>
                <title>{props.pageTitle}</title>
                <link rel="icon" href="/favicon.ico"/>
                <meta name="viewport" content="width=device-width, initial-scale=1"/>
                <meta charSet="utf-8"/>

                <meta name="description" content={props.description}/>
                <meta property="og:title" content={props.pageTitle} key="ogtitle"/>
                <meta property="og:description" content={props.description} key="ogdesc"/>

                {/* Open Graph */}
                {/*<meta property="og:url" content={props.currentURL} key="ogurl" />*/}
                <meta property="og:image" content={props.previewImage} key="ogimage"/>
                <meta property="og:site_name" content={props.pageTitle} key="ogsitename"/>

            </Head>

           {/* REMOVE HEADER */}
            <header className="bg-white dark:bg-gray-800">
            <IzeNavbar />
            
                <div className="lg:flex">
                    <div className="flex justify-center w-full px-6 py-6 lg:h-128 lg:w-1/2">
                        <div className="max-w-xl">
                            <TypeItAnimation />

                            <p className="mt-20 text-sm text-gray-500 dark:text-gray-400 lg:text-base">
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
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                              d="M5 13l4 4L19 7"/>
                                    </svg>

                                    <span>AWS</span>
                                </div>

                                <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                    <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                         viewBox="0 0 24 24"
                                         stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                              d="M5 13l4 4L19 7"/>
                                    </svg>

                                    <span>Terraform</span>
                                </div>

                                <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                    <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                         viewBox="0 0 24 24"
                                         stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                              d="M5 13l4 4L19 7"/>
                                    </svg>

                                    <span>Docker</span>
                                </div>

                                <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                    <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                         viewBox="0 0 24 24"
                                         stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                              d="M5 13l4 4L19 7"/>
                                    </svg>

                                    <span>Serverless</span>
                                </div>

                                <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                    <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                         viewBox="0 0 24 24"
                                         stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                              d="M5 13l4 4L19 7"/>
                                    </svg>

                                    <span>ECS</span>
                                </div>

                                <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
                                    <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none"
                                         viewBox="0 0 24 24"
                                         stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                              d="M5 13l4 4L19 7"/>
                                    </svg>

                                    <span>CI/CD</span>
                                </div>
                            </div>



                        </div>


                    </div>


                    <div className="flex justify-center w-full px-6 py-8 lg:h-256 lg:w-1/2">
                        <div className="w-full py-12">
                            <div className="coding inverse-toggle px-5 pt-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased
                 bg-gray-800  pb-6 pt-4 rounded-lg leading-normal overflow-hidden text-left">
                                <div className="top mb-2 flex">
                                    <div className="h-3 w-3 bg-red-500 rounded-full"></div>
                                    <div className="ml-2 h-3 w-3 bg-yellow-500 rounded-full"></div>
                                    <div className="ml-2 h-3 w-3 bg-green-500 rounded-full"></div>
                                </div>

                                <div className="mt-4 flex">
                                    <p className="flex-1 typing items-center pl-0">
                                        <span className="text-gray-400"># Deploy infrastructure</span><br/>
                                        <span className="text-blue-600">❯</span> ize deploy infra<br/>
                                    </p>
                                </div>

                                <div className="mt-4 flex">
                                    <p className="flex-1 typing items-center pl-0">
                                        <span className="text-gray-400"># Build a web service</span><br/>
                                        <span className="text-blue-600">❯</span> ize build web<br/>
                                    </p>
                                </div>


                                <div className="mt-4 flex">
                                    <p className="flex-1 typing items-center pl-0">
                                        <span className="text-gray-400"># Deploy a web service</span><br/>
                                        <span className="text-blue-600">❯</span> ize deploy web<br/>
                                    </p>
                                </div>


                                <div className="mt-4 flex">
                                    <p className="flex-1 typing items-center pl-0">
                                        <span className="text-gray-400"># Connect to SSH Bastion Tunnel</span><br/>
                                        <span className="text-blue-600">❯</span> ize tunnel up<br/>
                                    </p>
                                </div>


                                <div className="mt-4 flex">
                                    <p className="flex-1 typing items-center pl-0">
                                        <span className="text-gray-400"># Connect to the web service container via SSH/SSM</span><br/>
                                        <span className="text-blue-600">❯</span> ize ssh web<br/>
                                    </p>
                                </div>

                            </div>
                        </div>
                    </div>

                </div>


            </header>
            <section className="text-gray-600 body-font">
                <div className="container px-5 py-24 mx-auto">
                    <div className="flex flex-wrap w-full mb-20 flex-col items-center text-center">
                        <h1 id="features" className="sm:text-3xl text-2xl font-medium title-font mb-2 text-gray-900">Features</h1>
                    </div>
                    <div className="flex flex-wrap -m-4">
                        <div className="xl:w-1/3 md:w-1/2 p-4">
                            <div className="border border-gray-200 p-6 rounded-lg">
                                <div
                                    className="w-8 h-8 text-blue-600">
                                    <FontAwesomeIcon icon="fa-solid fa-diagram-project" />
                                </div>
                                <h2 className="text-lg text-gray-900 font-medium title-font mb-2">Coherent Infrastructure Deployment</h2>
                                <p className="leading-relaxed text-base">We abstract infrastructure management and
                                    provide a clean coherent way to deploy it.</p>
                                <p className="leading-relaxed text-base">We integrate with the following tools to
                                    perform infra rollouts:</p>
                                <p className="leading-relaxed text-base">- Terraform</p>
                                <p className="leading-relaxed text-base">- Ansible*</p>
                                <p className="leading-relaxed text-base">- Cloudformation*</p>
                                {/*<p className="leading-relaxed text-base">*Currently under development</p>*/}
                            </div>
                        </div>
                        <div className="xl:w-1/3 md:w-1/2 p-4">
                            <div className="border border-gray-200 p-6 rounded-lg">
                                <div
                                    className="w-8 h-8 text-blue-600">
                                    <FontAwesomeIcon icon="fa-solid fa-list-check" />
                                </div>
                                <h2 className="text-lg text-gray-900 font-medium title-font mb-2">Coherent Application Deployment</h2>
                                <p className="leading-relaxed text-base">We unify application deployment process and
                                    utilize naming conventions to facilitate and streamline deployments.</p>
                                <p className="leading-relaxed text-base">We allow to describe:</p>
                                <p className="leading-relaxed text-base">- ECS (currently using ecs-deploy
                                    underneath)</p>
                                <p className="leading-relaxed text-base">- k8s*</p>
                                <p className="leading-relaxed text-base">- Serverless*</p>
                                {/*<p className="leading-relaxed text-base">*Currently under development</p>*/}
                            </div>
                        </div>
                        <div className="xl:w-1/3 md:w-1/2 p-4">
                            <div className="border border-gray-200 p-6 rounded-lg">
                                <div
                                    className="w-8 h-8 text-blue-600">
                                    <FontAwesomeIcon icon="fa-solid fa-building-shield" />
                                </div>
                                <h2 className="text-lg text-gray-900 font-medium title-font mb-2">Port Forwarding via Bastion Host</h2>
                                <p className="leading-relaxed text-base">You don’t need to setup VPN solutions to you
                                    private network, if you are just starting out.</p>
                                <p className="leading-relaxed text-base"> Also you don’t need to compromise with
                                    security.</p>
                                <p className="leading-relaxed text-base">Establish port forwarding seamlessly to any
                                    private resource via your bastion host and connect to your private resources securely.</p>
                                {/*<p className="leading-relaxed text-base">*</p>*/}
                            </div>
                        </div>
                        <div className="xl:w-1/3 md:w-1/2 p-4">
                            <div className="border border-gray-200 p-6 rounded-lg">
                                <div
                                    className="w-8 h-8 text-blue-600">
                                    <FontAwesomeIcon icon="fa-solid fa-terminal" />
                                </div>
                                <h2 className="text-lg text-gray-900 font-medium title-font mb-2">Interactive Console to Fargate Containers</h2>
                                <p className="leading-relaxed text-base">You can access your containers running on AWS
                                    Fargate by providing the service name.
                                    </p>

                            </div>
                        </div>

                        <div className="xl:w-1/3 md:w-1/2 p-4">
                            <div className="border border-gray-200 p-6 rounded-lg">
                                <div
                                    className="w-8 h-8 text-blue-600">
                                    <FontAwesomeIcon icon="fa-solid fa-key" />
                                </div>
                                <h2 className="text-lg text-gray-900 font-medium title-font mb-2">Application Secrets Management</h2>
                                <p className="leading-relaxed text-base">Push, Remove your secrets to/from AWS Parameter
                                    Store.</p>

                            </div>
                        </div>
                        <div className="xl:w-1/3 md:w-1/2 p-4">
                            <div className="border border-gray-200 p-6 rounded-lg">
                                <div
                                    className="w-8 h-8 text-blue-600">
                                    <FontAwesomeIcon icon="fa-solid fa-seedling" />
                                </div>
                                <h2 className="text-lg text-gray-900 font-medium title-font mb-2">Terraform Environment Management</h2>
                                <p className="leading-relaxed text-base">Definitions of the environment can be stored in
                                    a toml file in the local repository.</p>
                            </div>
                        </div>
                    </div>
                    <h1 className="mt-16 italic border-0 py-2 px-8">*Currently under development</h1>
                </div>
            </section>

            <footer className="bg-white dark:bg-gray-800">
                <div className="container px-6 py-8 mx-auto">


                    <hr className="my-10 dark:border-gray-700"></hr>

                    <div className="sm:flex sm:items-center sm:justify-between">
                        <p className="text-sm text-gray-400">Ize is an Open Source software licensed under <a
                            href="https://raw.githubusercontent.com/hazelops/ize/main/LICENSE" target="_blank">Apache
                            2.0</a>. © Copyright 2022 <a href="https://hazelops.com" target="_blank">HazelOps OÜ</a>.
                        </p>

                        <div className="flex mt-3 -mx-2 sm:mt-0">


                            <a href="https://github.com/hazelops/ize"
                               className="mx-2 text-gray-400 hover:text-gray-500 dark:hover:text-gray-300"
                               aria-label="Github">
                                <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24" fill="none"
                                     xmlns="http://www.w3.org/2000/svg">
                                    <path
                                        d="M12.026 2C7.13295 1.99937 2.96183 5.54799 2.17842 10.3779C1.395 15.2079 4.23061 19.893 8.87302 21.439C9.37302 21.529 9.55202 21.222 9.55202 20.958C9.55202 20.721 9.54402 20.093 9.54102 19.258C6.76602 19.858 6.18002 17.92 6.18002 17.92C5.99733 17.317 5.60459 16.7993 5.07302 16.461C4.17302 15.842 5.14202 15.856 5.14202 15.856C5.78269 15.9438 6.34657 16.3235 6.66902 16.884C6.94195 17.3803 7.40177 17.747 7.94632 17.9026C8.49087 18.0583 9.07503 17.99 9.56902 17.713C9.61544 17.207 9.84055 16.7341 10.204 16.379C7.99002 16.128 5.66202 15.272 5.66202 11.449C5.64973 10.4602 6.01691 9.5043 6.68802 8.778C6.38437 7.91731 6.42013 6.97325 6.78802 6.138C6.78802 6.138 7.62502 5.869 9.53002 7.159C11.1639 6.71101 12.8882 6.71101 14.522 7.159C16.428 5.868 17.264 6.138 17.264 6.138C17.6336 6.97286 17.6694 7.91757 17.364 8.778C18.0376 9.50423 18.4045 10.4626 18.388 11.453C18.388 15.286 16.058 16.128 13.836 16.375C14.3153 16.8651 14.5612 17.5373 14.511 18.221C14.511 19.555 14.499 20.631 14.499 20.958C14.499 21.225 14.677 21.535 15.186 21.437C19.8265 19.8884 22.6591 15.203 21.874 10.3743C21.089 5.54565 16.9181 1.99888 12.026 2Z">
                                    </path>
                                </svg>
                            </a>
                        </div>
                    </div>
                </div>
            </footer>
        </div>
    )
}
