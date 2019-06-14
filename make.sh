#!/bin/bash
rm -f bin/src
#rm -f *.log
#$1为第一个参数
msg=$1
git add .
git commit -m" ${msg}"
git push


