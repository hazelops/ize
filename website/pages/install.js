// https://tailblocks.cc/
// https://merakiui.com/
// https://www.tailwind-kit.com/components#elements
import Head from 'next/head'
import IzeNavbar from '../components/izeNavbar'
// import TypeIt from "typeit-react"
import {Helmet} from 'react-helmet';

export default function Home() {
    let props = {
        pageTitle: "ize: Opinionated Infra Tool",
        description: "Opinionated Infra Tool",
        previewImage: "/social-preview.png"
    }

    return (

        <div className="flex flex-col">
            <Helmet>
                <link rel="stylesheet" href="install.css"/>
            </Helmet>
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

            <header className="bg-white dark:bg-gray-800">
            <IzeNavbar />
                <div class="container flex flex-wrap px-4 mx-auto items-center">
                    <div class="md:pr-12 md:py-8 mb-10 md:mb-0 pb-10 ">
                        <h1 id="installation"
                            class="sm:text-3xl text-2xl font-medium title-font mb-2 text-gray-900">Installation</h1>
                        <p class="leading-relaxed text-base"> Ize can be installed on MacOS and Ubuntu</p>
                        <h3>1. Install the latest version via homebrew on MacOS:</h3>
                        <h5 class="leading-relaxed text-base">1.1 Install <a target="_blank"
                                                                             href="https://brew.sh/">Homebrew</a></h5>
                        <h5 class="leading-relaxed text-base">1.2 Run the following commands:</h5>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> brew tap hazelops/ize <br/>
                                <span className="text-blue-600">❯</span> brew install ize <br/>
                            </p>
                        </div>

                        <p class="leading-relaxed text-base"> Now you can run ize from command shell by
                            typing <code>ize</code> in console.</p>

                        <h5 class="text-xl font-medium mb-2 text-gray-700"> 2. Update ize:</h5>

                        <p class="leading-relaxed text-base"> To update ize you should:</p>


                        <h5 class="leading-relaxed text-base">2.1 Uninstall previous version:</h5>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> brew uninstall ize <br/>
                            </p>
                        </div>

                        <h5 class="leading-relaxed text-base">2.2 Update version in brew repo:</h5>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> brew tap hazelops/ize <br/>
                            </p>
                        </div>

                        <h5 class="leading-relaxed text-base">2.3 Install:</h5>


                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> brew install ize <br/>
                            </p>
                        </div>

                        <h3 class="sm:text-xl text-xl font-medium mb-2 text-gray-700">3. Installation via public apt
                            repository URL (Ubuntu):</h3>


                        <h5 class="leading-relaxed text-base">3.1 To add public apt repository run:</h5>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span>echo "deb [trusted=yes]
                                https://apt.fury.io/hazelops/ /"|sudo tee /etc/apt/sources.list.d/fury.list<br/>
                            </p>
                        </div>

                        <h5 class="leading-relaxed text-base">3.2 After this, you should update information. Run:</h5>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> sudo apt-get update <br/>
                            </p>
                        </div>


                        <h5 class="leading-relaxed text-base">3.3 To install the latest version of ize app, you should
                            run:</h5>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> sudo apt-get install ize <br/>
                            </p>
                        </div>

                        <h5 class="leading-relaxed text-base">3.4 If you wish to install certain version of the ize you
                            should add version like this:</h5>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> sudo apt-get install ize=ize_version <br/>
                            </p>
                        </div>

                        <h5 class="leading-relaxed text-base">3.5 To remove ize app - run this command:</h5>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> sudo apt-get purge ize <br/>
                            </p>
                        </div>


                        <h3 class="sm:text-xl text-xl font-medium mb-2 text-gray-700">4. Installation from source:</h3>

                        <ul>
                            <li>GO version should be 1.16+</li>
                            <li><code>GOPATH</code> environment variable is set to <code>~/go</code></li>
                        </ul>

                        <p class="leading-relaxed text-base">To install ize from source download code or clone it from
                            this repo. After this you should run:</p>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> go mod download <br/>
                                <span className="text-blue-600">❯</span> mod install <br/>
                            </p>
                        </div>


                        <h3 class="sm:text-xl text-xl font-medium mb-2 text-gray-700">5. Autocomplete:</h3>

                        <p class="leading-relaxed text-base"> You could use integrated option to add autocompletion to
                            ize commands (bash, fish, zsh, powershell). In this manual we will describe it only for zsh
                            and bash. </p>


                        <p class="leading-relaxed text-base">To add autocompletion script, use the following manual:</p>

                        <h5 class="leading-relaxed text-base">5.1 ZSH:</h5>

                        <p class="leading-relaxed text-base">If shell completion is not already enabled in your
                            environment you will need to enable it. You should execute the following once:</p>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> echo "autoload -U compinit; compinit" >>
                                ~/.zshrc<br/>
                            </p>
                        </div>

                        <p class="leading-relaxed text-base">To load completions for every new session, execute
                            once:</p>

                        <h5 class="leading-relaxed text-base">5.1.1 macOS:</h5>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> ize completion zsh >
                                /usr/local/share/zsh/site-functions/_ize <br/>
                            </p>
                        </div>

                        <h5 class="leading-relaxed text-base">5.1.2 Linux:</h5>

                        <p class="leading-relaxed text-base">You will need root privileges.</p>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> sudo zsh <br/>
                            </p>
                        </div>


                        <p class="leading-relaxed text-base">Input your root password and run:</p>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">


                                <span className="text-blue-600">❯</span> soft completion zsh > {'"${fpath[0]}"'} <br/>
                            </p>
                        </div>


                        <p class="leading-relaxed text-base">To take effect for this setup you should run source
                            ~/.zshrc or simply restart shell.</p>


                        <h5 class="leading-relaxed text-base">5.2 Bash:</h5>

                        <p class="leading-relaxed text-base">Autocompletion script depends on the bash-completion
                            package. If it is not installed already, you can install it via your OS package manager.</p>

                        <h5 class="leading-relaxed text-base">5.2.1 MacOS:</h5>

                        <p class="leading-relaxed text-base">To load completions for every new session, you should
                            execute once:</p>

                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> ize completion bash >
                                /usr/local/etc/bash_completion.d/ize <br/>
                            </p>
                        </div>


                        <h5 class="leading-relaxed text-base">5.2.2 Linux:</h5>

                        <p class="leading-relaxed text-base">You will need root privileges.</p>


                        <div
                            className="lg:w-1/1 coding inverse-toggle px-5 pt-4 mb-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased bg-gray-800  rounded-lg leading-normal overflow-hidden text-left">

                            <p className="flex-1 typing items-center pl-0">
                                <span className="text-blue-600">❯</span> sudo bash <br/>
                                <span className="text-blue-600">❯</span> ize completion bash >
                                /etc/bash_completion.d/ize <br/>
                            </p>
                        </div>


                    </div>
                    <p class="leading-relaxed text-base">To take effect for this setup you should run source ~/.bashrc
                        or simply restart shell.</p>
                </div>
            </header>


            <footer className="bg-white dark:bg-gray-800">
                <div className="container px-6 py-8 mx-auto">


                    <hr className="my-10 dark:border-gray-700"></hr>

                    <div className="sm:flex sm:items-center sm:justify-between">
                        <p className="text-sm text-gray-400">Ize is an Open Source software licensed under <a
                            href="https://raw.githubusercontent.com/hazelops/ize/main/LICENSE" target="_blank">Apache
                            2.0</a>. © Copyright 2021 <a href="https://hazelops.com" target="_blank">HazelOps OÜ</a>.
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
