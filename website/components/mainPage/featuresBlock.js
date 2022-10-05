import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

function FeatureLayout({ title, icon, children }) {
    return (
        // xl:w-1/3 md:w-1/2 sm:w-1/2
        <div className="p-4 xl:w-1/3 md:w-1/2 sm:w-full w-full">
            <div className="border border-gray-200 px-6 pt-6 rounded-lg h-full">
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

function Paragraph({ text }) {
    return <p className="leading-relaxed text-base pb-3">{text}</p>
}

function Feature({ feature }) {
    const { icon, title, content } = feature
    let renderContent

    if (typeof content === "object" && content.length === undefined) {
        const paragraphs = Object.keys(content)
        renderContent = paragraphs.map((el, ind) => {
            if (content[el]) {
                const list = content[el].map((listEl, ind) => {
                    return (
                        <li key={ind} className="leading-relaxed text-base list-inside">
                            <span className="pl-1">{listEl}</span>
                        </li>
                    )
                })
                return (
                    <div key={ind} className="pb-3">
                        <Paragraph text={el} />
                        <ul className="list-['-']">
                            {list}
                        </ul>
                    </div>
                )
            }
            return <Paragraph key={ind} text={el} />
        })
    } else if (typeof content === "object" && content.length != undefined) {
        renderContent = content.map((el, ind) => {
            return <Paragraph key={ind} text={el} />
        })
    } else {
        renderContent = <Paragraph text={content} />
    }

    return (
        <FeatureLayout title={title} icon={icon}>
            {renderContent}
        </FeatureLayout>
    )
}

// ---------------------------------------------------------------

export default function FeaturesBlock({ extraData, features }) {
    const { header, underDev } = extraData
    const listFeatures = features.map((feature, ind) => {
        return <Feature key={ind} feature={feature} />
    })

    return (
        <section className="text-gray-600 body-font">
            <div className="container px-5 py-24 mx-auto">
                <div className="flex flex-wrap w-full mb-10 flex-col items-center text-center">
                    <h1 id="features" className="sm:text-3xl text-2xl font-medium title-font mb-2 text-gray-900">{header}</h1>
                </div>

                <div className="flex flex-wrap justify-between">
                    {listFeatures}
                </div>

                <h1 className="mt-16 italic border-0 py-2 px-8">{underDev}</h1>
            </div>
        </section>
    )
}
