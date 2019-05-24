FROM node:lts-alpine as vuebuild
WORKDIR /
COPY client/package*.json /
RUN npm install
COPY client/ /
RUN npm run build

FROM golang:alpine as gobuild
WORKDIR /go/src/github.com/shreve/fileserver
COPY server .
RUN go build -o files .

FROM nginx:alpine
WORKDIR /root/
VOLUME /files
ENV FILES_ROOT /files
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=vuebuild /dist /usr/share/nginx/html
COPY --from=gobuild /go/src/github.com/shreve/fileserver/files .
RUN ln -s /usr/share/nginx/html/index.html ./index.html
COPY ./run.sh .
ENTRYPOINT ["sh", "run.sh"]
