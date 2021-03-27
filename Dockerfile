FROM node:14 as frontend

COPY ./frontend /opt/talkstimer
WORKDIR /opt/talkstimer

RUN rm -f src/api/const.js
RUN mv src/api/const.js.prod src/api/const.js

RUN npm i
RUN npm run build

FROM golang:alpine as backend

COPY ./server /go/src/github.com/meyskens/TalksTimer/server
WORKDIR /go/src/github.com/meyskens/TalksTimer/server


RUN CGO_ENABLED=0 go build -o talkstimer ./

FROM alpine:3.13

RUN apk add --no-cache ca-certificates
COPY --from=backend /go/src/github.com/meyskens/TalksTimer/server/talkstimer /opt/talkstimer/talkstimer
COPY --from=frontend /opt/talkstimer/ /opt/talkstimer/www/

WORKDIR /opt/talkstimer
CMD /opt/talkstimer/talkstimer
