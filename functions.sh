#!/bin/bash
#include this file use: source functions.sh

#usage: result=`inArray $search "${list[*]}"`
function inArray() {
  #Confine variables inside functions using the local command
  local search=$1
  local array=$2
  local result=0
  for item in ${array[*]}; do
    if [[ $item == $search ]]; then
      result=1
    fi
  done
  echo $result
}
