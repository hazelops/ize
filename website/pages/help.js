// https://tailblocks.cc/
// https://merakiui.com/
// https://www.tailwind-kit.com/components#elements
import Head from 'next/head'

export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen py-2">
      <Head>
        <title>Create Next Apps</title>
        <link rel="icon" href="/favicon.ico"/>
      </Head>

      <main className="flex flex-col items-center justify-center w-full flex-1 px-20 text-center">

        <div className="max-w-7xl">

          <div className="coding inverse-toggle px-5 pt-4 shadow-lg text-gray-100 text-sm font-mono subpixel-antialiased
              bg-gray-800  pb-6 pt-4 rounded-lg leading-normal overflow-hidden text-left">
            <div className="top mb-2 flex">
              <div className="h-3 w-3 bg-red-500 rounded-full"></div>
              <div className="ml-2 h-3 w-3 bg-yellow-500 rounded-full"></div>
              <div className="ml-2 h-3 w-3 bg-green-500 rounded-full"></div>
            </div>

            <div className="mt-4 flex">

              <span className="text-green-400">$</span>

              <p className="flex-1 typing items-center pl-2">
                brew tap hazelops/ize<br/>
              </p>
            </div>
            <div className="mt-4 flex">

              <span className="text-green-400">$</span>

              <p className="flex-1 typing items-center pl-2">
                brew install ize<br/>
              </p>
            </div>
          </div>
        </div>

        {/*<section className="text-gray-600 body-font ">*/}
        {/*  <div className="container px-5 py-24 mx-auto ">*/}
        {/*    <div className="flex flex-col text-center w-full mb-20">*/}
        {/*      <h2 className="text-xs text-indigo-500 tracking-widest font-medium title-font mb-1">ROOF PARTY POLAROID</h2>*/}
        {/*      <h1 className="sm:text-3xl text-2xl font-medium title-font text-gray-900">Master Cleanse Reliac Heirloom</h1>*/}
        {/*    </div>*/}
        {/*    <div className="flex flex-wrap -m-4">*/}
        {/*      <div className="p-4 md:w-1/3">*/}
        {/*        <div className="flex rounded-lg h-full bg-gray-100 p-8 flex-col">*/}
        {/*          <div className="flex items-center mb-3">*/}
        {/*            <div className="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-indigo-500 text-white flex-shrink-0">*/}
        {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-5 h-5" viewBox="0 0 24 24">*/}
        {/*                <path d="M22 12h-4l-3 9L9 3l-3 9H2"></path>*/}
        {/*              </svg>*/}
        {/*            </div>*/}
        {/*            <h2 className="text-gray-900 text-lg title-font font-medium">Shooting Stars</h2>*/}
        {/*          </div>*/}
        {/*          <div className="flex-grow">*/}
        {/*            <p className="leading-relaxed text-base">Blue bottle crucifix vinyl post-ironic four dollar toast vegan taxidermy. Gastropub indxgo juice poutine.</p>*/}
        {/*            <a className="mt-3 text-indigo-500 inline-flex items-center">Learn More*/}
        {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-4 h-4 ml-2" viewBox="0 0 24 24">*/}
        {/*                <path d="M5 12h14M12 5l7 7-7 7"></path>*/}
        {/*              </svg>*/}
        {/*            </a>*/}
        {/*          </div>*/}
        {/*        </div>*/}
        {/*      </div>*/}
        {/*      <div class="p-4 md:w-1/3">*/}
        {/*        <div class="flex rounded-lg h-full bg-gray-100 p-8 flex-col">*/}
        {/*          <div class="flex items-center mb-3">*/}
        {/*            <div class="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-indigo-500 text-white flex-shrink-0">*/}
        {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-5 h-5" viewBox="0 0 24 24">*/}
        {/*                <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"></path>*/}
        {/*                <circle cx="12" cy="7" r="4"></circle>*/}
        {/*              </svg>*/}
        {/*            </div>*/}
        {/*            <h2 class="text-gray-900 text-lg title-font font-medium">The Catalyzer</h2>*/}
        {/*          </div>*/}
        {/*          <div class="flex-grow">*/}
        {/*            <p class="leading-relaxed text-base">Blue bottle crucifix vinyl post-ironic four dollar toast vegan taxidermy. Gastropub indxgo juice poutine.</p>*/}
        {/*            <a class="mt-3 text-indigo-500 inline-flex items-center">Learn More*/}
        {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-4 h-4 ml-2" viewBox="0 0 24 24">*/}
        {/*                <path d="M5 12h14M12 5l7 7-7 7"></path>*/}
        {/*              </svg>*/}
        {/*            </a>*/}
        {/*          </div>*/}
        {/*        </div>*/}
        {/*      </div>*/}
        {/*      <div class="p-4 md:w-1/3">*/}
        {/*        <div class="flex rounded-lg h-full bg-gray-100 p-8 flex-col">*/}
        {/*          <div class="flex items-center mb-3">*/}
        {/*            <div class="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-indigo-500 text-white flex-shrink-0">*/}
        {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-5 h-5" viewBox="0 0 24 24">*/}
        {/*                <circle cx="6" cy="6" r="3"></circle>*/}
        {/*                <circle cx="6" cy="18" r="3"></circle>*/}
        {/*                <path d="M20 4L8.12 15.88M14.47 14.48L20 20M8.12 8.12L12 12"></path>*/}
        {/*              </svg>*/}
        {/*            </div>*/}
        {/*            <h2 class="text-gray-900 text-lg title-font font-medium">Neptune</h2>*/}
        {/*          </div>*/}
        {/*          <div class="flex-grow">*/}
        {/*            <p class="leading-relaxed text-base">Blue bottle crucifix vinyl post-ironic four dollar toast vegan taxidermy. Gastropub indxgo juice poutine.</p>*/}
        {/*            <a class="mt-3 text-indigo-500 inline-flex items-center">Learn More*/}
        {/*              <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="w-4 h-4 ml-2" viewBox="0 0 24 24">*/}
        {/*                <path d="M5 12h14M12 5l7 7-7 7"></path>*/}
        {/*              </svg>*/}
        {/*            </a>*/}
        {/*          </div>*/}
        {/*        </div>*/}
        {/*      </div>*/}
        {/*    </div>*/}
        {/*  </div>*/}
        {/*</section>*/}
      </main>

      {/*<footer className="flex items-center justify-center w-full h-24 border-t">*/}
      {/*  /!*<a*!/*/}
      {/*  /!*  className="flex items-center justify-center"*!/*/}
      {/*  /!*  href="https://vercel.com?utm_source=create-next-app&utm_medium=default-template&utm_campaign=create-next-app"*!/*/}
      {/*  /!*  target="_blank"*!/*/}
      {/*  /!*  rel="noopener noreferrer"*!/*/}
      {/*  /!*>*!/*/}
      {/*  /!*  Powered by{' '}*!/*/}
      {/*  /!*  <img src="/vercel.svg" alt="Vercel Logo" className="h-4 ml-2"/>*!/*/}
      {/*  /!*</a>*!/*/}
      {/*</footer>*/}
    </div>
  )
}
