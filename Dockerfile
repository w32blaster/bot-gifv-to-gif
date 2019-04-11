#
# Stage 1: Install dependencies
#
FROM alpine AS builder-ffmpeg

ARG FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz


# Download ffmpeg
ADD ${FFMPEG_URL} /tmp/ffmpeg.tar.xz
RUN cd /tmp && tar xJf ffmpeg.tar.xz


#
# Stage 2: compile and install giflossy
#
FROM alpine:3.8 AS builder-giflossy

RUN apk add --no-cache curl autoconf automake make build-base
RUN curl -SL https://github.com/kornelski/giflossy/archive/1.91.tar.gz | tar xzv
RUN cd giflossy-1.91 && autoreconf -i && ./configure && make install
RUN cp "$(which gifsicle)" /tmp/gifsicle


#
# Stage 3: Compile and build Go app
#

FROM golang:alpine AS builder-go

# Install Git for go get
RUN set -eux; \
    apk add --no-cache --virtual git

ENV GO_WORKDIR $GOPATH/src/github.com/w32blaster/bot-gifv-to-gif/
WORKDIR $GO_WORKDIR

ADD . $GO_WORKDIR
RUN go get -u gopkg.in/telegram-bot-api.v4 && \
    go get -u github.com/gofrs/uuid

# RUN TESTS HERE AS WELL!

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bot .
RUN cp bot /tmp

#
# Stage 4: build final image
#
FROM scratch

COPY --from=builder-ffmpeg /tmp/ffmpeg*/ffmpeg /bin/
COPY --from=builder-giflossy /tmp/gifsicle /bin/
COPY --from=builder-go /tmp/bot /bin/

# copy root CA certificate to set up HTTPS connection with Telegram
COPY --from=builder-go /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/bin/bot"]
