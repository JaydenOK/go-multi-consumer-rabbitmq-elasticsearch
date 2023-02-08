#!/bin/bash
source functions.sh

function start() {
  pId=$(ps aux | grep $targetFile | grep -v grep | awk '{print $2}')
  if [[ $pId != "" ]]; then
    echo "kill exist server, pid: ${pid}"
    kill -9 $pId
  fi
  # run server in backstage daemon
  nohup $binDir/$targetFile -env ${env} >>${logFile} 2>&1 &
  echo "server started"
}

function stop() {
  pId=$(ps aux | grep $targetFile | grep -v grep | awk '{print $2}')
  if [[ $pId != "" ]]; then
    echo "kill server, pid: ${pid}"
    kill -9 $pId
    echo "done"
  else
    echo "server not running"
  fi
}

function status() {
  pId=$(ps aux | grep $targetFile | grep -v grep | awk '{print $2}')
  if [[ $pId != "" ]]; then
    echo "server running, pid: ${pId}"
  else
    echo "server not running"
  fi
}

function list() {
  msg=$(ps aux | grep $outputSuffix | grep -v grep)
  if [[ $msg == "" ]]; then
    echo "no server running"
  else
    echo $msg
  fi
}

# usage:
# ./manage.sh stop myapp dev
# ./manage.sh start myapp dev
# ./manage.sh status myapp dev
# ./manage.sh list

scriptDir=$(
  cd $(dirname $0)
  pwd
)

appRoot=$scriptDir
binDir=${appRoot}/bin
logRir=${appRoot}/logs
outputSuffix=go_app

cmd=$1
appName=$2
env=$3

if [[ $cmd == "" || $cmd != "list" && ($appName == "" || $env == "") ]]; then
  echo 'usage: ./manage.sh cmd[start|stop|status] app[appName] env[dev|test|prod]'
  exit
fi

targetFile=${appName}_${env}_${outputSuffix}
logFile=${targetFile}_$(date '+%Y-%m').log

cmdList=("start" "stop" "status" "list")
if [[ $(inArray ${cmd} "${cmdList[*]}") != 1 ]]; then
  echo "cmd not exist:${cmd}"
  exit
fi

envList=("dev" "test" "prod")
if [[ $cmd != "list" && $(inArray ${env} "${envList[*]}") != 1 ]]; then
  echo "env not exist:${env}"
  exit
fi

case $cmd in
"start")
  start
  ;;
"stop")
  stop
  ;;
"status")
  status
  ;;
"list")
  list
  ;;
*)
  echo 'error cmd'
  ;;
esac
