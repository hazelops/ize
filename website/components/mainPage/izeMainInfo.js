function Feature({ title }) {
    return (
        <div className="flex items-center space-x-6 text-gray-800 dark:text-gray-200">
            <svg className="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
            </svg>

            <span>{title}</span>
        </div>
    )
}

export default function IzeMainInfo({ izeDescription, mainFeatures }) {
    const listFeatures = mainFeatures.map((title, ind) => {
        return <Feature key={ind} title={title} />
    })
    return (
        <>
            <p className="mt-6 text-sm text-gray-500 dark:text-gray-400 lg:text-base">
                {izeDescription}
            </p>

            <div className="grid gap-6 mt-8 sm:grid-cols-2">
                {listFeatures}
            </div>
        </>
    )
}
