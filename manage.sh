#!/bin/bash

scriptDir=$(
  cd $(dirname $0)
  pwd
)

function start() {
  pId=$(ps aux | grep $targetFile | grep -v grep | awk '{print $2}')
  if [[ $pId != "" ]]; then
    echo "kill exist process pid: {$pid}"
    kill -9 $pId
  fi
  # run server in backstage daemon
  nohup $targetFile -env ${env} >>${logFile} 2>&1 &
  echo "server started"
}

function stop() {
  pId=$(ps aux | grep $targetFile | grep -v grep | awk '{print $2}')
  if [[ $pId != "" ]]; then
    echo "kill exist server process, pid: {$pid}"
    kill -9 $pId
  else
    echo "server not running"
  fi
}

function status() {
  pId=$(ps aux | grep $targetFile | grep -v grep | awk '{print $2}')
  if [[ $pId != "" ]]; then
    echo "server running, pid: {$pId}"
  else
    echo "server not running"
  fi
}

# usage:
# ./manage.sh stop myapp dev
# ./manage.sh start myapp dev
# ./manage.sh status myapp dev

appRoot=$scriptDir
binDir=${appRoot}/bin
logRir=${appRoot}/logs

cmd=$1
appName=$2
env=$3

if [[ $cmd == "" || $appName == "" || $env == "" ]]; then
  echo 'command usage:./manage.sh cmd[start|stop|status] app[appName] env[dev|test|prod]'
  exit
fi

targetFile=$logRir/${appName}_go_app
logFile=${targetFile}_$(date '+%Y-%m').log

envList=("dev","test","prod")
if [[ "${envList[@]}" =~ "${env}" ]]; then
   echo ""
else
    echo "env not config:{$env}"
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
*)
  echo 'error'
  ;;
esac
