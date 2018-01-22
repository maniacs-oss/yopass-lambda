#!/bin/bash
if [ "$S3_BUCKET" == "" ]; then
echo "specify S3_BUCKET environment variable"
exit 1
fi

(cd create;GOOS=linux go build -o ../createHandler;)
(cd get;GOOS=linux go build -o ../getHandler;)

zip deployment.zip main createHandler
aws cloudformation package --template-file template.yml --s3-bucket "$S3_BUCKET" --output-template-file packaged.yml
aws cloudformation deploy --template-file /Users/jhaals/go/src/github.com/yopass/yopass-lambda/packaged.yml --stack-name yopass --capabilities CAPABILITY_IAM