export const data = {
    header: "Installation via public apt repository URL (Ubuntu):",
    content: {
        "To add public apt repository run:": [ 
            'echo "deb [trusted=yes] https://apt.fury.io/hazelops/ /"|sudo tee /etc/apt/sources.list.d/fury.list'
        ],
        "After this, you should update information. Run:": [ 
            "sudo apt-get update" 
        ],
        "To install the latest version of ize app, you should run:": [ 
            "sudo apt-get install ize"
        ],
        "If you wish to install certain version of the ize you should add version like this:": [ 
            "sudo apt-get install ize=ize_version"
        ],
        "To remove ize app - run this command:": [ "sudo apt-get purge ize" ]
    }
}
