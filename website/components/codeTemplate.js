export default function CodeTemplate({ children }) {
    return (
<div className="mt-5 w-1/2 flex items-center justify-center">
        <div className="mx-5 w-1/2 bg-gray-800 rounded-lg">
            <div id="header-buttons" className="py-3 px-4 flex">
                <div className="rounded-full w-3 h-3 bg-red-500 mr-2"></div>
                <div className="rounded-full w-3 h-3 bg-yellow-500 mr-2"></div>
                <div className="rounded-full w-3 h-3 bg-green-500"></div>
            </div>
            <div id="code-area" className="py-4 px-4 mt-1 text-white text-xl">
                {children}
            </div>
        </div>
    </div>
    )
}
