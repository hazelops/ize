import Ize from './ize'
import Link from 'next/link'
import React from 'react'
import GitHubIcon from './gitHubIcon'

export default function IzeNavbar() {
    return (
            <nav className="bg-white shadow dark:bg-gray-800" >
                    <div className="container inline-block px-6 py-4 mx-auto lg:flex lg:justify-between lg:items-center">
                        <div className="lg:flex lg:items-center">
                            <div className="flex justify-between">
                                <Ize />

                                {/*Mobile Button*/}
                                <div className="flex lg:hidden">
                                    <button type="button"
                                            className="text-gray-500 dark:text-gray-200 hover:text-gray-600 dark:hover:text-gray-400 focus:outline-none focus:text-gray-600 dark:focus:text-gray-400"
                                            aria-label="toggle menu">
                                        <svg viewBox="0 0 24 24" className="w-6 h-6 fill-current">
                                            <path fillRule="evenodd"
                                                  d="M4 5h16a1 1 0 0 1 0 2H4a1 1 0 1 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2z"></path>
                                        </svg>
                                    </button>
                                </div>

                                <div className="flex flex-col text-gray-600 capitalize dark:text-gray-300 lg:flex lg:px-16 lg:-mx-4 lg:flex-row lg:items-center">
                                    <Link href="/docs/installation">
                                        <a className="mt-2 lg:mt-0 lg:mx-4 hover:text-gray-800 dark:hover:text-gray-200">
                                            installation
                                        </a>
                                    </Link>  

                                    <Link href="/docs">
                                        <a className="mt-2 lg:mt-0 lg:mx-4 hover:text-gray-800 dark:hover:text-gray-200">
                                            docs
                                        </a>
                                    </Link>
                                    
                                    <a href="#features" className="mt-2 lg:mt-0 lg:mx-4 hover:text-gray-800 dark:hover:text-gray-200">
                                        features
                                    </a>
                                </div>
                            </div>   
                        </div>

                        <GitHubIcon />
                    </div>
            </nav>
    )
}
