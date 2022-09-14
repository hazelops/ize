import SideBar from '../components/sideBar'
import DocBody from '../components/docBody'

export default function DocsLayout({ data }) {
    return (
        <div className="flex">
            <SideBar />
            <div>
                <DocBody data={data}/>
            </div> 
        </div>
    )
}
