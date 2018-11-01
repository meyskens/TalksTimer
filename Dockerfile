ARG ARCH
FROM node:10 as frontend

COPY ./ /opt/talkstimer
WORKDIR /opt/talkstimer

RUN npm i
RUN npm run build

FROM golang as backend

COPY ./ /go/src/github.com/meyskens/TalksTimer/server
WORKDIR /go/src/github.com/meyskens/TalksTimer/server

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} GOARM=${GOARM} go build -a -installsuffix cgo -o talkstimer ./

# Set up deinitive image
ARG ARCH
FROM multiarch/alpine:${ARCH}-edge

RUN apk add --no-cache ca-certificates
COPY --from=backend /go/src/github.com/meyskens/TalksTimer/talkstimer /opt/talkstimer/talkstimer
COPY --from=frontend /opt/talkstimer/dist/ /opt/talkstimer/www/

CMD talkstimer
