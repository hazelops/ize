export const extraData = {
    header: "Features",
    underDev: "*Currently under development"
}

export const features = [
    {
        icon: "fa-solid fa-diagram-project",
        title: "Coherent Infrastructure Deployment",
        content: {
            "We abstract infrastructure management and provide a clean coherent way to deploy it.": null,
            "We integrate with the following tools to perform infra rollouts:": [
                "Terraform",
                "Ansible*",
                "Cloudformation*"
            ]
        }
    },
    {
        icon: "fa-solid fa-list-check",
        title: "Coherent Application Deployment",
        content: {
           "We unify application deployment process and utilize naming conventions to facilitate and streamline deployments.": null,
            "We allow to describe:": [
                "ECS (currently using ecs-deploy underneath)",
                "k8s*",
                "Serverless*"
            ] 
        }  
    },
    {
        icon: "fa-solid fa-building-shield",
        title: "Port Forwarding via Bastion Host",
        content: [
            "You don’t need to setup VPN solutions to you private network, if you are just starting out.",
            "Also you don’t need to compromise with security.",
            "Establish port forwarding seamlessly to any private resource via your bastion host and connect to your private resources securely."
        ] 
    },
    {
        icon: "fa-solid fa-terminal",
        title: "Interactive Console to Fargate Containers",
        content: "You can access your containers running on AWS Fargate by providing the service name."
    }, 
    {
        icon: "fa-solid fa-key",
        title: "Application Secrets Management",
        content: "Push, Remove your secrets to/from AWS Parameter Store."
    },
    {
        icon: "fa-solid fa-seedling",
        title: "Terraform Environment Management",
        content: "Definitions of the environment can be stored in a toml file in the local repository."
    }
]
