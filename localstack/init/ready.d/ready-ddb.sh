#!/bin/bash

awslocal dynamodb create-table \
    --table-name goto-dev-todo \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

awslocal dynamodb list-tables