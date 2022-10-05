export const data = {
    header: "Features",
    underDev: "*Currently under development"
}

export const features = {
    "Coherent Infrastructure Deployment": {
        "We abstract infrastructure management and provide a clean coherent way to deploy it.": null,
        "We integrate with the following tools to perform infra rollouts:": [
            "Terraform",
            "Ansible*",
            "Cloudformation*"
        ]
    },
    "Coherent Application Deployment": {
        "We unify application deployment process and utilize naming conventions to facilitate and streamline deployments.": null,
        "We allow to describe:": [
            "ECS (currently using ecs-deploy underneath)",
            "k8s*",
            "Serverless*"
        ]
    },
    "Port Forwarding via Bastion Host": [
        "You don’t need to setup VPN solutions to you private network, if you are just starting out.",
        "Also you don’t need to compromise with security.",
        "Establish port forwarding seamlessly to any private resource via your bastion host and connect to your private resources securely."
    ],
    "Interactive Console to Fargate Containers": "You can access your containers running on AWS Fargate by providing the service name.",
    "Application Secrets Management": "Push, Remove your secrets to/from AWS Parameter Store.",
    "Terraform Environment Management": "Definitions of the environment can be stored in a toml file in the local repository."
}
