export const data = {
    install: {
        header: "Install the latest version via homebrew on MacOS:",
        content: {
            "InstallHomebrew": null,
            "Run the following commands:": [
                "brew tap hazelops/ize",
                "brew install ize"
            ]
        },
        extraData: "Now you can run ize from command shell by typing ize in console."
    },
    update: {
        header: "Update ize:",
        extraData: "To update ize you should:",
        content: {
            "Uninstall previous version:": [ "brew uninstall ize" ],
            "Update version in brew repo:": [ "brew tap hazelops/ize" ],
            "Install:": [ "brew install ize" ]
        }
    }
}
