import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

function Feature({ props }) {
    return (
        null
    )
}

function FeaturBlock({ data, icon, children }) {
    return (
        <div className="xl:w-1/3 md:w-1/2 p-4">
            <div className="border border-gray-200 p-6 rounded-lg">
                <div
                    className="w-8 h-8 text-blue-600">
                    <FontAwesomeIcon icon={icon} />
                </div>
                
                <h2 className="text-lg text-gray-900 font-medium title-font mb-2">{title}</h2>
                {children}
            </div>
        </div>
    )
}

export default function FeaturesBlock({ extraData, features, children }) {
    const { header, underDev } = extraData
    
    return (
        <section className="text-gray-600 body-font">
            <div className="container px-5 py-24 mx-auto">
                <div className="flex flex-wrap w-full mb-20 flex-col items-center text-center">
                    <h1 id="features" className="sm:text-3xl text-2xl font-medium title-font mb-2 text-gray-900">{header}</h1>
                </div>

                <div className="flex flex-wrap -m-4">
                           {/* blocks */}
                </div>

                <h1 className="mt-16 italic border-0 py-2 px-8">{underDev}</h1>
            </div>
        </section>
    )
}
