#!/bin/bash

#Install aws cli


if which aws > /dev/null
    then
        echo "AWS CLI found."
    else
        echo "Installing AWS CLI"

        if ! [ -f awscli-bundle.zip ]
            then
                curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
        fi

        unzip -n awscli-bundle.zip
        sudo ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

        echo "Install complete; cleaning up bundle files."
        rm -r ./awscli-bundle

fi

clear

export AWS_ACCESS_KEY_ID=AKIAYGK2HOIZXVM6GOF4
export AWS_SECRET_ACCESS_KEY=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

cat <<EOF ~/.aws/config
[default]
output = json
region = us-west-2
EOF

cat <<EOF ~/.aws/credentials
[default]
aws_access_key_id = AKIAYGK2HOIZXVM6GOF4
aws_secret_access_key = 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
EOF

aws eks update-kubeconfig --name SRVAWSEKSCl01 --region=us-west-2

echo -e "Variable defined"

echo -e "\nGetting cluster"
kubectl get all --all-namespaces

#salvar como configure-eks.sh
#sudo chmod +x configure-eks.sh