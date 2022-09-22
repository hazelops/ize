import CodeTemplate from "../codeTemplate";

export default function MacOSInstallation() {
    return (
        <>
            <h3>
                1. Install the latest version via homebrew on MacOS:
            </h3>
            <h5 className="leading-relaxed text-base">
                1.1 Install
                <a target="_blank" href="https://brew.sh/">Homebrew</a>
            </h5>                    
            <h5 className="leading-relaxed text-base">
                1.2 Run the following commands:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span> 
                brew tap hazelops/ize <br/>
                <span className="text-blue-600">❯</span> 
                brew install ize <br/>
            </CodeTemplate>
            <p className="leading-relaxed text-base">
                Now you can run ize from command shell by typing <code>ize</code> in console.
            </p>
            <h5 className="text-xl font-medium mb-2 text-gray-700">
                2. Update ize:
            </h5>
            <p className="leading-relaxed text-base">
                To update ize you should:
            </p>
            <h5 className="leading-relaxed text-base">
                2.1 Uninstall previous version:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span> 
                brew uninstall ize <br/>
            </CodeTemplate>
            <h5 className="leading-relaxed text-base">
                2.2 Update version in brew repo:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span> 
                brew tap hazelops/ize <br/>
            </CodeTemplate>
            <h5 className="leading-relaxed text-base">
                2.3 Install:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span> 
                brew install ize <br/>
            </CodeTemplate>
        </>
    )
}
