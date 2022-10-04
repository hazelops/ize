export const data = {
    extraInfo: "You could use integrated option to add autocompletion to ize commands (bash, fish, zsh, powershell). In this manual we will describe it only for zsh and bash.",
    header: "To add autocompletion script, use the following manual:",
    content: {
        "ZSH:": {
            extraInfo: [
                "If shell completion is not already enabled in your environment you will need to enable it. You should execute the following once:",
                'echo "autoload -U compinit; compinit" >> ~/.zshrc',
                "To load completions for every new session, execute once:"
            ],
            content: {
                
            }
        },
        "Bash:": {

        }
    }
}
