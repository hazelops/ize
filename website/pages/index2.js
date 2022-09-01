// https://tailblocks.cc/
// https://merakiui.com/
// https://www.tailwind-kit.com/components#elements
import Head from 'next/head'

export default function Home() {
  return (


    // <nav className="bg-white shadow dark:bg-gray-800">
    //   <div className="container px-6 py-4 mx-auto lg:flex lg:justify-between lg:items-center">
    //     <div className="lg:flex lg:items-center">
    //       <div className="flex items-center justify-between">
    //         <div>
    //           <a
    //             className="text-2xl font-bold text-gray-800 dark:text-white lg:text-3xl hover:text-gray-700 dark:hover:text-gray-300"
    //             href="#">❯ IZE</a>
    //         </div>
    //         {/*<!-- Mobile menu button -->*/}
    //
    //         <div className="flex lg:hidden">
    //           <button type="button"
    //                   className="text-gray-500 dark:text-gray-200 hover:text-gray-600 dark:hover:text-gray-400 focus:outline-none focus:text-gray-600 dark:focus:text-gray-400"
    //                   aria-label="toggle menu">
    //             <svg viewBox="0 0 24 24" className="w-6 h-6 fill-current">
    //               <path fill-rule="evenodd"
    //                     d="M4 5h16a1 1 0 0 1 0 2H4a1 1 0 1 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2z"></path>
    //             </svg>
    //           </button>
    //         </div>
    //       </div>
    //
    //       <div
    //         className="flex flex-col text-gray-600 capitalize dark:text-gray-300 lg:flex lg:px-16 lg:-mx-4 lg:flex-row lg:items-center">
    //
    //         {/*<a href="#" className="mt-2 lg:mt-0 lg:mx-4 hover:text-gray-800 dark:hover:text-gray-200">features</a>*/}
    //         {/*<a href="#" className="mt-2 lg:mt-0 lg:mx-4 hover:text-gray-800 dark:hover:text-gray-200">downloads</a>*/}
    //         {/*<a href="#" className="mt-2 lg:mt-0 lg:mx-4 hover:text-gray-800 dark:hover:text-gray-200">docs</a>*/}
    //         {/*<a href="#" className="mt-2 lg:mt-0 lg:mx-4 hover:text-gray-800 dark:hover:text-gray-200">support</a>*/}
    //         {/*<a href="#" className="mt-2 lg:mt-0 lg:mx-4 hover:text-gray-800 dark:hover:text-gray-200">blog</a>*/}
    //
    //         {/*<div className="relative mt-4 lg:mt-0 lg:mx-4">*/}
    //         {/*          <span className="absolute inset-y-0 left-0 flex items-center pl-3">*/}
    //         {/*              <svg className="w-4 h-4 text-gray-600 dark:text-gray-300" viewBox="0 0 24 24" fill="none">*/}
    //         {/*                  <path*/}
    //         {/*                    d="M21 21L15 15M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z"*/}
    //         {/*                    stroke="currentColor" stroke-width="2" stroke-linecap="round"*/}
    //         {/*                    stroke-linejoin="round"></path>*/}
    //         {/*              </svg>*/}
    //         {/*          </span>*/}
    //
    //         {/*  <input type="text"*/}
    //         {/*         className="w-full py-1 pl-10 pr-4 text-gray-700 placeholder-gray-600 bg-white border-b border-gray-600 dark:placeholder-gray-300 dark:focus:border-gray-300 lg:w-56 lg:border-transparent dark:bg-gray-800 dark:text-gray-300 focus:outline-none focus:border-gray-600"*/}
    //         {/*         placeholder="Search">*/}
    //         {/*</div>*/}
    //       </div>
    //     </div>
    //
    //     <div className="flex justify-center mt-6 lg:flex lg:mt-0 lg:-mx-2">
    //
    //
    //       {/*<a href="#" className="mx-2 text-gray-600 dark:text-gray-300 hover:text-gray-500 dark:hover:text-gray-300"*/}
    //       {/*   aria-label="Reddit">*/}
    //       {/*  <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24" fill="none"*/}
    //       {/*       xmlns="http://www.w3.org/2000/svg">*/}
    //       {/*    <path*/}
    //       {/*      d="M12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12C21.9939 17.5203 17.5203 21.9939 12 22ZM6.807 10.543C6.20862 10.5433 5.67102 10.9088 5.45054 11.465C5.23006 12.0213 5.37133 12.6558 5.807 13.066C5.92217 13.1751 6.05463 13.2643 6.199 13.33C6.18644 13.4761 6.18644 13.6229 6.199 13.769C6.199 16.009 8.814 17.831 12.028 17.831C15.242 17.831 17.858 16.009 17.858 13.769C17.8696 13.6229 17.8696 13.4761 17.858 13.33C18.4649 13.0351 18.786 12.3585 18.6305 11.7019C18.475 11.0453 17.8847 10.5844 17.21 10.593H17.157C16.7988 10.6062 16.458 10.7512 16.2 11C15.0625 10.2265 13.7252 9.79927 12.35 9.77L13 6.65L15.138 7.1C15.1931 7.60706 15.621 7.99141 16.131 7.992C16.1674 7.99196 16.2038 7.98995 16.24 7.986C16.7702 7.93278 17.1655 7.47314 17.1389 6.94094C17.1122 6.40873 16.6729 5.991 16.14 5.991C16.1022 5.99191 16.0645 5.99491 16.027 6C15.71 6.03367 15.4281 6.21641 15.268 6.492L12.82 6C12.7983 5.99535 12.7762 5.993 12.754 5.993C12.6094 5.99472 12.4851 6.09583 12.454 6.237L11.706 9.71C10.3138 9.7297 8.95795 10.157 7.806 10.939C7.53601 10.6839 7.17843 10.5422 6.807 10.543ZM12.18 16.524C12.124 16.524 12.067 16.524 12.011 16.524C11.955 16.524 11.898 16.524 11.842 16.524C11.0121 16.5208 10.2054 16.2497 9.542 15.751C9.49626 15.6958 9.47445 15.6246 9.4814 15.5533C9.48834 15.482 9.52348 15.4163 9.579 15.371C9.62737 15.3318 9.68771 15.3102 9.75 15.31C9.81233 15.31 9.87275 15.3315 9.921 15.371C10.4816 15.7818 11.159 16.0022 11.854 16C11.9027 16 11.9513 16 12 16C12.059 16 12.119 16 12.178 16C12.864 16.0011 13.5329 15.7863 14.09 15.386C14.1427 15.3322 14.2147 15.302 14.29 15.302C14.3653 15.302 14.4373 15.3322 14.49 15.386C14.5985 15.4981 14.5962 15.6767 14.485 15.786V15.746C13.8213 16.2481 13.0123 16.5208 12.18 16.523V16.524ZM14.307 14.08H14.291L14.299 14.041C13.8591 14.011 13.4994 13.6789 13.4343 13.2429C13.3691 12.8068 13.6162 12.3842 14.028 12.2269C14.4399 12.0697 14.9058 12.2202 15.1478 12.5887C15.3899 12.9572 15.3429 13.4445 15.035 13.76C14.856 13.9554 14.6059 14.0707 14.341 14.08H14.306H14.307ZM9.67 14C9.11772 14 8.67 13.5523 8.67 13C8.67 12.4477 9.11772 12 9.67 12C10.2223 12 10.67 12.4477 10.67 13C10.67 13.5523 10.2223 14 9.67 14Z">*/}
    //       {/*    </path>*/}
    //       {/*  </svg>*/}
    //       {/*</a>*/}
    //
    //       {/*<a href="#" className="mx-2 text-gray-600 dark:text-gray-300 hover:text-gray-500 dark:hover:text-gray-300"*/}
    //       {/*   aria-label="Facebook">*/}
    //       {/*  <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24" fill="none"*/}
    //       {/*       xmlns="http://www.w3.org/2000/svg">*/}
    //       {/*    <path*/}
    //       {/*      d="M2.00195 12.002C2.00312 16.9214 5.58036 21.1101 10.439 21.881V14.892H7.90195V12.002H10.442V9.80204C10.3284 8.75958 10.6845 7.72064 11.4136 6.96698C12.1427 6.21332 13.1693 5.82306 14.215 5.90204C14.9655 5.91417 15.7141 5.98101 16.455 6.10205V8.56104H15.191C14.7558 8.50405 14.3183 8.64777 14.0017 8.95171C13.6851 9.25566 13.5237 9.68693 13.563 10.124V12.002H16.334L15.891 14.893H13.563V21.881C18.8174 21.0506 22.502 16.2518 21.9475 10.9611C21.3929 5.67041 16.7932 1.73997 11.4808 2.01722C6.16831 2.29447 2.0028 6.68235 2.00195 12.002Z">*/}
    //       {/*    </path>*/}
    //       {/*  </svg>*/}
    //       {/*</a>*/}
    //
    //       <a href="https://github.com/hazelops/ize" className="mx-2 text-gray-600 dark:text-gray-300 hover:text-gray-500 dark:hover:text-gray-300"
    //          aria-label="Github">
    //         <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24" fill="none"
    //              xmlns="http://www.w3.org/2000/svg">
    //           <path
    //             d="M12.026 2C7.13295 1.99937 2.96183 5.54799 2.17842 10.3779C1.395 15.2079 4.23061 19.893 8.87302 21.439C9.37302 21.529 9.55202 21.222 9.55202 20.958C9.55202 20.721 9.54402 20.093 9.54102 19.258C6.76602 19.858 6.18002 17.92 6.18002 17.92C5.99733 17.317 5.60459 16.7993 5.07302 16.461C4.17302 15.842 5.14202 15.856 5.14202 15.856C5.78269 15.9438 6.34657 16.3235 6.66902 16.884C6.94195 17.3803 7.40177 17.747 7.94632 17.9026C8.49087 18.0583 9.07503 17.99 9.56902 17.713C9.61544 17.207 9.84055 16.7341 10.204 16.379C7.99002 16.128 5.66202 15.272 5.66202 11.449C5.64973 10.4602 6.01691 9.5043 6.68802 8.778C6.38437 7.91731 6.42013 6.97325 6.78802 6.138C6.78802 6.138 7.62502 5.869 9.53002 7.159C11.1639 6.71101 12.8882 6.71101 14.522 7.159C16.428 5.868 17.264 6.138 17.264 6.138C17.6336 6.97286 17.6694 7.91757 17.364 8.778C18.0376 9.50423 18.4045 10.4626 18.388 11.453C18.388 15.286 16.058 16.128 13.836 16.375C14.3153 16.8651 14.5612 17.5373 14.511 18.221C14.511 19.555 14.499 20.631 14.499 20.958C14.499 21.225 14.677 21.535 15.186 21.437C19.8265 19.8884 22.6591 15.203 21.874 10.3743C21.089 5.54565 16.9181 1.99888 12.026 2Z">
    //           </path>
    //         </svg>
    //       </a>
    //     </div>
    //   </div>
    // </nav>

    // <div className="flex flex-col items-center justify-center min-h-screen py-2">
    //   <Head>
    //     <title>❯ ize: Opinionated Infra Tool</title>
    //     <link rel="icon" href="/favicon.ico"/>
    //   </Head>
    //
    //
    //
    //
    //
    //   <main className="flex flex-col items-center justify-center w-full flex-1 px-20 text-center">
    //
    //     <div className="max-w-7xl">
    //
    //       <div className="coding inverse-toggle px-5 pt-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased
    //           bg-gray-800  pb-6 pt-4 rounded-lg leading-normal overflow-hidden text-left">
    //         <div className="top mb-2 flex">
    //           <div className="h-3 w-3 bg-red-500 rounded-full"></div>
    //           <div className="ml-2 h-3 w-3 bg-yellow-500 rounded-full"></div>
    //           <div className="ml-2 h-3 w-3 bg-green-500 rounded-full"></div>
    //         </div>
    //
    //         <div className="mt-4 flex">
    //
    //           <span className="text-green-400">$</span>
    //
    //           <p className="flex-1 typing items-center pl-2">
    //             brew tap hazelops/ize<br/>
    //           </p>
    //         </div>
    //         <div className="mt-4 flex">
    //
    //           <span className="text-green-400">$</span>
    //
    //           <p className="flex-1 typing items-center pl-2">
    //             brew install ize<br/>
    //           </p>
    //         </div>
    //       </div>
    //     </div>
    //
    //     {/*<section className="text-gray-600 body-font ">*/}
    //     {/*  <div className="container px-5 py-24 mx-auto ">*/}
    //     {/*    <div className="flex flex-col text-center w-full mb-20">*/}
    //     {/*      <h2 className="text-xs text-indigo-500 tracking-widest font-medium title-font mb-1">ROOF PARTY POLAROID</h2>*/}
    //     {/*      <h1 className="sm:text-3xl text-2xl font-medium title-font text-gray-900">Master Cleanse Reliac Heirloom</h1>*/}
    //     {/*    </div>*/}
    //     {/*    <div className="flex flex-wrap -m-4">*/}
    //     {/*      <div className="p-4 md:w-1/3">*/}
    //     {/*        <div className="flex rounded-lg h-full bg-gray-100 p-8 flex-col">*/}
    //     {/*          <div className="flex items-center mb-3">*/}
    //     {/*            <div className="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-indigo-500 text-white flex-shrink-0">*/}
    //     {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-5 h-5" viewBox="0 0 24 24">*/}
    //     {/*                <path d="M22 12h-4l-3 9L9 3l-3 9H2"></path>*/}
    //     {/*              </svg>*/}
    //     {/*            </div>*/}
    //     {/*            <h2 className="text-gray-900 text-lg title-font font-medium">Shooting Stars</h2>*/}
    //     {/*          </div>*/}
    //     {/*          <div className="flex-grow">*/}
    //     {/*            <p className="leading-relaxed text-base">Blue bottle crucifix vinyl post-ironic four dollar toast vegan taxidermy. Gastropub indxgo juice poutine.</p>*/}
    //     {/*            <a className="mt-3 text-indigo-500 inline-flex items-center">Learn More*/}
    //     {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-4 h-4 ml-2" viewBox="0 0 24 24">*/}
    //     {/*                <path d="M5 12h14M12 5l7 7-7 7"></path>*/}
    //     {/*              </svg>*/}
    //     {/*            </a>*/}
    //     {/*          </div>*/}
    //     {/*        </div>*/}
    //     {/*      </div>*/}
    //     {/*      <div class="p-4 md:w-1/3">*/}
    //     {/*        <div class="flex rounded-lg h-full bg-gray-100 p-8 flex-col">*/}
    //     {/*          <div class="flex items-center mb-3">*/}
    //     {/*            <div class="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-indigo-500 text-white flex-shrink-0">*/}
    //     {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-5 h-5" viewBox="0 0 24 24">*/}
    //     {/*                <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"></path>*/}
    //     {/*                <circle cx="12" cy="7" r="4"></circle>*/}
    //     {/*              </svg>*/}
    //     {/*            </div>*/}
    //     {/*            <h2 class="text-gray-900 text-lg title-font font-medium">The Catalyzer</h2>*/}
    //     {/*          </div>*/}
    //     {/*          <div class="flex-grow">*/}
    //     {/*            <p class="leading-relaxed text-base">Blue bottle crucifix vinyl post-ironic four dollar toast vegan taxidermy. Gastropub indxgo juice poutine.</p>*/}
    //     {/*            <a class="mt-3 text-indigo-500 inline-flex items-center">Learn More*/}
    //     {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-4 h-4 ml-2" viewBox="0 0 24 24">*/}
    //     {/*                <path d="M5 12h14M12 5l7 7-7 7"></path>*/}
    //     {/*              </svg>*/}
    //     {/*            </a>*/}
    //     {/*          </div>*/}
    //     {/*        </div>*/}
    //     {/*      </div>*/}
    //     {/*      <div class="p-4 md:w-1/3">*/}
    //     {/*        <div class="flex rounded-lg h-full bg-gray-100 p-8 flex-col">*/}
    //     {/*          <div class="flex items-center mb-3">*/}
    //     {/*            <div class="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-indigo-500 text-white flex-shrink-0">*/}
    //     {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-5 h-5" viewBox="0 0 24 24">*/}
    //     {/*                <circle cx="6" cy="6" r="3"></circle>*/}
    //     {/*                <circle cx="6" cy="18" r="3"></circle>*/}
    //     {/*                <path d="M20 4L8.12 15.88M14.47 14.48L20 20M8.12 8.12L12 12"></path>*/}
    //     {/*              </svg>*/}
    //     {/*            </div>*/}
    //     {/*            <h2 class="text-gray-900 text-lg title-font font-medium">Neptune</h2>*/}
    //     {/*          </div>*/}
    //     {/*          <div class="flex-grow">*/}
    //     {/*            <p class="leading-relaxed text-base">Blue bottle crucifix vinyl post-ironic four dollar toast vegan taxidermy. Gastropub indxgo juice poutine.</p>*/}
    //     {/*            <a class="mt-3 text-indigo-500 inline-flex items-center">Learn More*/}
    //     {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-4 h-4 ml-2" viewBox="0 0 24 24">*/}
    //     {/*                <path d="M5 12h14M12 5l7 7-7 7"></path>*/}
    //     {/*              </svg>*/}
    //     {/*            </a>*/}
    //     {/*          </div>*/}
    //     {/*        </div>*/}
    //     {/*      </div>*/}
    //     {/*    </div>*/}
    //     {/*  </div>*/}
    //     {/*</section>*/}
    //   </main>
    //
    //   {/*<footer className="flex items-center justify-center w-full h-24 border-t">*/}
    //   {/*  /!*<a*!/*/}
    //   {/*  /!*  className="flex items-center justify-center"*!/*/}
    //   {/*  /!*  href="https://vercel.com?utm_source=create-next-app&utm_medium=default-template&utm_campaign=create-next-app"*!/*/}
    //   {/*  /!*  target="_blank"*!/*/}
    //   {/*  /!*  rel="noopener noreferrer"*!/*/}
    //   {/*  /!*>*!/*/}
    //   {/*  /!*  Powered by{' '}*!/*/}
    //   {/*  /!*  <img src="/vercel.svg" alt="Vercel Logo" className="h-4 ml-2"/>*!/*/}
    //   {/*  /!*</a>*!/*/}
    //   {/*</footer>*/}
    // </div>


    <header className="bg-white dark:bg-gray-800">
      <nav className="bg-white dark:bg-gray-800">
        <div className="container p-6 mx-auto">
          <a
            className="block text-2xl font-bold text-center text-gray-800 dark:text-white lg:text-3xl hover:text-gray-700 dark:hover:text-gray-300"
            href="#">Brand</a>

          <div className="flex items-center justify-center mt-6 text-gray-600 capitalize dark:text-gray-300">
            <a href="#" className="text-gray-800 dark:text-gray-200 border-b-2 border-blue-500 mx-1.5 sm:mx-6">home</a>

            <a href="#"
               className="border-b-2 border-transparent hover:text-gray-800 dark:hover:text-gray-200 hover:border-blue-500 mx-1.5 sm:mx-6">features</a>

            <a href="#"
               className="border-b-2 border-transparent hover:text-gray-800 dark:hover:text-gray-200 hover:border-blue-500 mx-1.5 sm:mx-6">pricing</a>

            <a href="#"
               className="border-b-2 border-transparent hover:text-gray-800 dark:hover:text-gray-200 hover:border-blue-500 mx-1.5 sm:mx-6">blog</a>

            <a href="#"
               className="border-b-2 border-transparent hover:text-gray-800 dark:hover:text-gray-200 hover:border-blue-500 mx-1.5 sm:mx-6">
              <svg className="w-4 h-4 fill-current" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" clip-rule="evenodd"
                      d="M1 11.9554V12.0446C1.01066 14.7301 1.98363 17.1885 3.59196 19.0931C4.05715 19.6439 4.57549 20.1485 5.13908 20.5987C5.70631 21.0519 6.31937 21.4501 6.97019 21.7853C7.90271 22.2656 8.91275 22.6165 9.97659 22.8143C10.5914 22.9286 11.2243 22.9918 11.8705 22.9993C11.9136 22.9998 11.9567 23 11.9999 23C15.6894 23 18.9547 21.1836 20.9502 18.3962C21.3681 17.8125 21.7303 17.1861 22.0291 16.525C22.6528 15.1448 22.9999 13.613 22.9999 12C22.9999 8.73978 21.5816 5.81084 19.3283 3.79653C18.8064 3.32998 18.2397 2.91249 17.6355 2.55132C15.9873 1.56615 14.0597 1 11.9999 1C11.888 1 11.7764 1.00167 11.6653 1.00499C9.99846 1.05479 8.42477 1.47541 7.0239 2.18719C6.07085 2.67144 5.19779 3.29045 4.42982 4.01914C3.7166 4.69587 3.09401 5.4672 2.58216 6.31302C2.22108 6.90969 1.91511 7.54343 1.6713 8.20718C1.24184 9.37631 1.00523 10.6386 1 11.9554ZM20.4812 15.0186C20.8171 14.075 20.9999 13.0588 20.9999 12C20.9999 9.54265 20.0151 7.31533 18.4186 5.6912C17.5975 7.05399 16.5148 8.18424 15.2668 9.0469C15.7351 10.2626 15.9886 11.5603 16.0045 12.8778C16.7692 13.0484 17.5274 13.304 18.2669 13.6488C19.0741 14.0252 19.8141 14.487 20.4812 15.0186ZM15.8413 14.8954C16.3752 15.0321 16.904 15.22 17.4217 15.4614C18.222 15.8346 18.9417 16.3105 19.5723 16.8661C18.0688 19.2008 15.5151 20.7953 12.5788 20.9817C13.5517 20.0585 14.3709 18.9405 14.972 17.6514C15.3909 16.7531 15.678 15.8272 15.8413 14.8954ZM13.9964 12.6219C13.9583 11.7382 13.7898 10.8684 13.5013 10.0408C10.6887 11.2998 7.36584 11.3765 4.35382 9.97197C4.01251 9.81281 3.68319 9.63837 3.36632 9.44983C3.12787 10.2584 2.99991 11.1142 2.99991 12C2.99991 13.9462 3.61763 15.748 4.6677 17.2203C6.83038 14.1875 10.3685 12.4987 13.9964 12.6219ZM6.047 18.7502C7.77258 16.059 10.7714 14.5382 13.8585 14.6191C13.723 15.3586 13.4919 16.093 13.1594 16.8062C12.3777 18.4825 11.1453 19.805 9.67385 20.6965C8.31043 20.3328 7.07441 19.6569 6.047 18.7502ZM11.9999 3C13.7846 3 15.4479 3.51946 16.847 4.41543C16.2113 5.54838 15.3593 6.4961 14.368 7.23057C13.3472 5.57072 11.8752 4.16433 10.027 3.21692C10.6619 3.07492 11.3222 3 11.9999 3ZM8.80619 4.84582C10.4462 5.61056 11.7474 6.80659 12.6379 8.23588C10.3464 9.24654 7.64722 9.30095 5.19906 8.15936C4.83384 7.98905 4.48541 7.79735 4.15458 7.58645C4.91365 6.24006 6.00929 5.10867 7.32734 4.30645C7.82672 4.44058 8.32138 4.61975 8.80619 4.84582Z"
                      fill="currentColor"/>
              </svg>
            </a>

            <a href="#"
               className="border-b-2 border-transparent hover:text-gray-800 dark:hover:text-gray-200 hover:border-blue-500 mx-1.5 sm:mx-6">
              <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path
                  d="m.75 19h7.092c4.552 0 6.131-6.037 2.107-8.203 2.701-2.354 1.029-6.797-2.595-6.797h-6.604c-.414 0-.75.336-.75.75v13.5c0 .414.336.75.75.75zm.75-13.5h5.854c3.211 0 3.215 4.768 0 4.768h-5.854zm0 6.268h6.342c3.861 0 3.861 5.732 0 5.732h-6.342z"/>
                <path
                  d="m18.374 7.857c-3.259 0-5.755 2.888-5.635 5.159-.247 3.28 2.397 5.984 5.635 5.984 2.012 0 3.888-1.065 4.895-2.781.503-.857-.791-1.613-1.293-.76-.739 1.259-2.12 2.041-3.602 2.041-2.187 0-3.965-1.668-4.125-3.771 1.443.017 4.136-.188 8.987-.033.016 0 .027-.008.042-.008 2-.09-.189-5.831-4.904-5.831zm-3.928 4.298c1.286-3.789 6.718-3.676 7.89.064-4.064.097-6.496-.066-7.89-.064z"/>
                <path d="m21.308 6.464c.993 0 .992-1.5 0-1.5h-5.87c-.993 0-.992 1.5 0 1.5z"/>
              </svg>
            </a>
          </div>
        </div>
      </nav>

      <div
        className="container flex flex-col px-6 py-4 mx-auto space-y-6 lg:h-128 lg:py-16 lg:flex-row lg:items-center lg:space-x-6">
        <div className="flex flex-col items-center w-full lg:flex-row lg:w-1/2">
          <div className="flex justify-center order-2 mt-6 lg:mt-0 lg:space-y-3 lg:flex-col">
            <button className="w-3 h-3 mx-2 bg-blue-500 rounded-full lg:mx-0 focus:outline-none"></button>
            <button
              className="w-3 h-3 mx-2 bg-gray-300 rounded-full lg:mx-0 focus:outline-none hover:bg-blue-500"></button>
            <button
              className="w-3 h-3 mx-2 bg-gray-300 rounded-full lg:mx-0 focus:outline-none hover:bg-blue-500"></button>
            <button
              className="w-3 h-3 mx-2 bg-gray-300 rounded-full lg:mx-0 focus:outline-none hover:bg-blue-500"></button>
          </div>

          <div className="max-w-lg lg:mx-12 lg:order-2">
            <h1 className="text-3xl font-medium tracking-wide text-gray-800 dark:text-white lg:text-4xl">The best Apple
              Watch apps</h1>
            <p className="mt-4 text-gray-600 dark:text-gray-300">Lorem ipsum, dolor sit amet consectetur adipisicing
              elit. Aut quia asperiores alias vero magnam recusandae adipisci ad vitae laudantium quod rem voluptatem
              eos accusantium cumque.</p>
            <div className="mt-6">
              <a href="#"
                 className="block px-3 py-2 font-semibold text-center text-white transition-colors duration-200 transform bg-blue-500 rounded-md lg:inline hover:bg-blue-400">Download
                from App Store</a>
            </div>
          </div>
        </div>

        <div className="flex items-center justify-center w-full h-96 lg:w-1/2">
          <img className="object-cover w-full h-full max-w-2xl rounded-md"
               src="https://images.unsplash.com/photo-1579586337278-3befd40fd17a?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=1352&q=80"
               alt="apple watch photo"></img>
        </div>
      </div>
    </header>



  )
}
