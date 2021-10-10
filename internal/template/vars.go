package template

var varsTemplate = `env="{{ .ENV }}"
aws_profile="{{ .AWS_PROFILE }}"
aws_region="{{ .AWS_REGION }}"
ec2_key_pair_name="{{ .EC2_KEY_PAIR_NAME }}"
docker_image_tag="{{ .TAG }}"
ssh_public_key="{{ .SSH_PUBLIC_KEY }}"
docker_registry="{{ .DOCKER_REGISTRY }}"
namespace="{{ .NAMESPACE }}"`
