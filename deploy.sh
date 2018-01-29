#!/bin/bash
GOOS=linux go build -o main

if [ "$S3_BUCKET" == "" ]; then
echo "specify S3_BUCKET environment variable"
exit 1
fi

zip deployment.zip main
aws cloudformation package --template-file template.yml --s3-bucket "$S3_BUCKET" --output-template-file packaged.yml
aws cloudformation deploy --template-file packaged.yml --stack-name yopass --capabilities CAPABILITY_IAM