import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

export default function FeatureBlock({ title, icon, children }) {
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
