import SideBar from '../components/sideBar'
import DocBody from '../components/docBody'

export default function DocsLayout({ data, filesNames }) {
    return (
        <div className="flex">
            <SideBar filesNames={filesNames} />
            <div>
                <DocBody data={data} />
            </div>
        </div>
    )
}
