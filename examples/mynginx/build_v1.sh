#!/bin/bash

set -euo pipefail

SITE="www.ntv.de"
# SITE="www.rtl.de"

wget -l1 --no-check-certificate -O www/index.html ${SITE}

MD5=$(md5 -q www/index.html|awk '{print substr($1,1,7)}')

docker build -t embano1/mynginx:${MD5} .
docker tag embano1/mynginx:${MD5} embano1/mynginx:v1
docker push embano1/mynginx:v1