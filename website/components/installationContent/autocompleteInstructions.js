import CodeTemplate from "../codeTemplate";

export default function AutocompleteInstructions() {
    return (
        <>
            <p class="leading-relaxed text-base">
                You could use integrated option to add autocompletion to ize commands (bash, fish, zsh, powershell).
                In this manual we will describe it only for zsh and bash. 
            </p>
            <p class="leading-relaxed text-base">
                To add autocompletion script, use the following manual:
            </p>
            <h5 class="leading-relaxed text-base">
                5.1 ZSH:
            </h5>
            <p class="leading-relaxed text-base">
            If shell completion is not already enabled in your environment you will need to enable it.
            You should execute the following once:
            </p>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                echo "autoload -U compinit; compinit" {'>'}{'>'} ~/.zshrc<br/>
            </CodeTemplate>
            <p class="leading-relaxed text-base">
                To load completions for every new session, execute once:
            </p>
            <h5 class="leading-relaxed text-base">
                5.1.1 macOS:
            </h5>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                ize completion zsh {'>'} /usr/local/share/zsh/site-functions/_ize <br/>
            </CodeTemplate>
            <h5 class="leading-relaxed text-base">
                5.1.2 Linux:
            </h5>
            <p class="leading-relaxed text-base">
                You will need root privileges.
            </p>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                sudo zsh <br/>
            </CodeTemplate>
            <p class="leading-relaxed text-base">
                Input your root password and run:
            </p>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                soft completion zsh {'>'} {'"${fpath[0]}"'} <br/>
            </CodeTemplate>
            <p class="leading-relaxed text-base">
                To take effect for this setup you should run source ~/.zshrc or simply restart shell.
            </p>
            <h5 class="leading-relaxed text-base">
                5.2 Bash:
            </h5>
            <p class="leading-relaxed text-base">
                Autocompletion script depends on the bash-completion package.
                If it is not installed already, you can install it via your OS package manager.
            </p>
            <h5 class="leading-relaxed text-base">
                5.2.1 MacOS:
            </h5>
            <p class="leading-relaxed text-base">
                To load completions for every new session, you should execute once:
            </p>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                ize completion bash {'>'} /usr/local/etc/bash_completion.d/ize <br/>
            </CodeTemplate>
            <h5 class="leading-relaxed text-base">
                5.2.2 Linux:
            </h5>
            <p class="leading-relaxed text-base">
                You will need root privileges.
            </p>
            <CodeTemplate>
                <span className="text-blue-600">❯</span>
                sudo bash <br/>
                <span className="text-blue-600">❯</span> 
                ize completion bash {'>'} /etc/bash_completion.d/ize <br/>
            </CodeTemplate>
            <p class="leading-relaxed text-base">
                To take effect for this setup you should run <code>source ~/.bashrc</code> or simply restart shell.
            </p>
        </> 
    )
}
