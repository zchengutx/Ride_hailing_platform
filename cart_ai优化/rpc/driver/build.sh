#!/bin/bash
RUN_NAME="driver"
mkdir -p output/bin
cp script/* output/
chmod +x output/bootstrap.sh
go build -o output/bin/${RUN_NAME} .
mkdir -p output/conf
cp ../basic/dev.yaml output/conf/
echo "司机服务构建完成" 