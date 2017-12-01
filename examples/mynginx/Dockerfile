FROM nginx:alpine
RUN apk add --no-cache wget
WORKDIR /tmp
# tag mynginx:v1
# ARG SITE="www.ntv.de"
# tag mynginx:v2
# ARG SITE="www.rtl.de"
# RUN wget -l1 --no-check-certificate -O index.html ${SITE} && cp index.html /usr/share/nginx/html
COPY www/index.html /usr/share/nginx/html