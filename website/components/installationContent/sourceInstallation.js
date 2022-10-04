import CodeTemplate from "../codeTemplate";
import { data } from "../../utilities/installationPage/source";

export default function SourceInstallation() {
    return (
        <>
            <h3 class="sm:text-xl text-xl font-medium mb-2 text-gray-700">4. Installation from source:</h3>
            <ul>
                <li>GO version should be 1.16+</li>
                <li><code>GOPATH</code> environment variable is set to <code>~/go</code></li>
            </ul>
            <p class="leading-relaxed text-base">
                To install ize from source download code or clone it from this repo. After this you should run:
            </p>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                go mod download <br/>
                <span className="text-blue-600">❯</span>
                mod install <br/>
            </CodeTemplate>
        </>
    )
}
