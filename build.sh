#!/bin/bash

scriptDir=$(cd `dirname $0`; pwd)

appRoot=$scriptDir
binDir=${appRoot}/bin
logRir=${appRoot}/logs

# get appName
while [[ $appName == "" ]]; do
    read -p "input appName:" appName
done

envList=("dev","test","prod")

# get env
while [[ $env == "" ]]; do
    read -p "input run env[dev|test|prod]:" env
    if [[ "${envList[@]}" =~ "${env}" ]];then
        break
    fi
done

echo -e "\n check again>  appName:${appName}, env:${env}"

targetFile=$logRir/${appName}_go_app

echo "go build ..."

#go build [-o outputName] [-i] [complimeMark] [packageName]
#go build -o /root/go/jayden/bin/app app

msg=$(go build -o $targetFile app 2>&1)
if [[ $msg == "" ]]; then
    echo "build success"
else
    echo "build fail:${msg}"
    exit
fi

# is start to run server
while [[ $isStart == "" ]]; do
    read -p "start compile and run server?[y/n]" isStart
    case $isStart in
        "y"|"Y")
            break
            ;;
        *)
            exit
            break
            ;;
    esac
done

pId=`ps aux|grep $targetFile|grep -v grep|awk '{print $2}'`
if [[ $pId != "" ]];then
    echo "kill exist server process pid: {$pId}"
    kill -9 $pId
fi

# run server in backstage daemon
logFile=${targetFile}_$(date '+%Y-%m').log
nohup $targetFile -env ${env} >> ${logFile} 2>&1 &
echo "server started"



#标记 描述
#-o 指定输出文件。
#-a 强行对所有涉及到的代码包（包括标准库中的代码包）进行重新构建，即使它们已经是最新的了。
#-n 打印构建期间所用到的其它命令，但是并不真正执行它们。
#-p n 构建的并行数量（n）。默认情况下并行数量与CPU数量相同。
#-race 开启数据竞争检测。此标记目前仅在linux/amd64、darwin/amd64和windows/amd64平台下被支持。
#-v 打印出被构建的代码包的名字。
#-work 打印出临时工作目录的名字，并且取消在构建完成后对它的删除操作。
#-x 打印出构建期间所用到的其它命令。