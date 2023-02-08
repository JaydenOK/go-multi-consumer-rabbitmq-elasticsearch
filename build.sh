#!/bin/bash

source functions.sh

outputSuffix=go_app
envList=("dev" "test" "prod")

scriptDir=$(
  cd $(dirname $0)
  pwd
)

appRoot=$scriptDir
binDir=${appRoot}/bin
logRir=${appRoot}/logs

# get appName
while [[ $appName == "" ]]; do
  read -p "input app name:" appName
done

# get env
while [[ $env == "" || $result == 0 ]]; do
  read -p "input environment[dev|test|prod]:" env
  result=$(inArray ${env} "${envList[*]}")
  if [[ $result == 1 ]]; then
    break
  fi
done

echo -e "\ncheck again > app name:${appName}, environment:${env}"

while [[ $confirm != "y" && $confirm != "Y" && $confirm != "n" && $confirm != "N" ]]; do
  read -p "continue build?[y/n]" confirm
  case $confirm in
  y | Y)
    break
    ;;
  n | N)
    exit
    break
    ;;
  esac
done

targetFile=${appName}_${env}_${outputSuffix}

echo "go build ..."

#go build [-o outputName] [-i] [complimeMark] [packageName]
#go build -o /root/go/jayden/bin/app app

msg=$(go build -o $binDir/$targetFile app 2>&1)
if [[ $msg == "" ]]; then
  echo "build success"
else
  echo "build fail:${msg}"
  exit
fi

# start server now
while [[ $isStart == "" || $isStart != "" && $isStart != "y" && $isStart != "Y" && $isStart != "n" && $isStart != "N" ]]; do
  read -p "start server now?[y/n]" isStart
  case $isStart in
  y | Y)
    break
    ;;
  n | N)
    exit
    break
    ;;
  esac
done

pId=$(ps aux | grep $targetFile | grep -v grep | awk '{print $2}')
if [[ $pId != "" ]]; then
  echo "kill exist server pid: {$pId}"
  kill -9 $pId
fi

# run server in backstage daemon
logFile=$logRir/${targetFile}_$(date '+%Y-%m').log
nohup $binDir/$targetFile -env ${env} >>${logFile} 2>&1 &
echo "server started"
