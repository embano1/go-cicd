#!/bin/bash

set -euo pipefail

SITE="www.ntv.de"
#SITE="www.rtl.de"
#SITE="www.stern.de"

wget -l1 --no-check-certificate -O www/index.html ${SITE}

MD5=$(md5 -q www/index.html|awk '{print substr($1,1,7)}')

docker build -t embano1/mynginx:${MD5} .
docker push embano1/mynginx:${MD5}
kubectl set image deploy/mynginx mynginx=embano1/mynginx:${MD5} --record
