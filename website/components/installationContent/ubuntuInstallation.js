import CodeTemplate from "../codeTemplate"
import { data } from "../../utilities/installationPage/ubuntu"

export default function UbuntuInstallation() {
    return (
        <>
            <h3 className="sm:text-xl text-xl font-medium mb-2 text-gray-700">
                3. Installation via public apt repository URL (Ubuntu):
            </h3>
            <h5 className="leading-relaxed text-base">
                3.1 To add public apt repository run:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                echo "deb [trusted=yes] https://apt.fury.io/hazelops/ /"|sudo tee /etc/apt/sources.list.d/fury.list <br/>
            </CodeTemplate>
            <h5 className="leading-relaxed text-base">
                3.2 After this, you should update information. Run:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                sudo apt-get update <br/>
            </CodeTemplate>
            <h5 className="leading-relaxed text-base">
                3.3 To install the latest version of ize app, you should run:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                sudo apt-get install ize <br/>
            </CodeTemplate>
            <h5 className="leading-relaxed text-base">
                3.4 If you wish to install certain version of the ize you should add version like this:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                sudo apt-get install ize=ize_version <br/>
            </CodeTemplate>
            <h5 className="leading-relaxed text-base">
                3.5 To remove ize app - run this command:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                sudo apt-get purge ize <br/>
            </CodeTemplate>
        </>
    )
}
