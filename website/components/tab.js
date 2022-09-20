function TabButton(props) {
    return (
        <button 
            class="h-10 px-4 py-2 -mb-px text-sm text-center text-blue-600 bg-transparent border-b-2 border-blue-500 sm:text-base whitespace-nowrap focus:outline-none"
            onClick={props.onClick}
        >
            {props.platform}
        </button>
    )
}

export default function Tab() {
    return (
        <div class="inline-flex border-b border-gray-200 dark:border-gray-700">
            <TabButton platform="MacOS" />
            <TabButton platform="Ubuntu" />
            <TabButton platform="Source" />
            <TabButton platform="Autocomplete" />
        </div>
    )
}
