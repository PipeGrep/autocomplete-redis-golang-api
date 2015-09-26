
FROM golang:1.5
MAINTAINER Samuel BERTHE <contact@samuel-berthe.fr>



# CONTAINER PARAMS

WORKDIR		/go
EXPOSE		8080
CMD		["bash", "-c", "service nginx start && make run"]







# NGINX

RUN             apt-key adv --keyserver pgp.mit.edu --recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62
RUN             echo "deb http://nginx.org/packages/mainline/debian/ wheezy nginx" >> /etc/apt/sources.list

RUN             apt-get update && apt-get install -y ca-certificates nginx
RUN             rm -rf /var/lib/apt/lists/*

COPY		nginx.conf /etc/nginx/nginx.conf






# GOLANG APP

ENV		GOROOT /usr/local/go
ENV		GOPATH /go/pkg:/go:/go/src

COPY		. /go/

RUN		cd /go && make vendor_clean && make vendor_get

