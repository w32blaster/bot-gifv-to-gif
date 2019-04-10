#
# Stage 1: Install dependencies
#
FROM alpine AS builder-ffmpeg

ARG FFMPEG_VERSION=4.1.3
ARG FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-$FFMPEG_VERSION-64bit-static.tar.xz

# Download ffmpeg
ADD ${FFMPEG_URL} /tmp/ffmpeg.tar.xz
RUN cd /tmp && tar xJf ffmpeg.tar.xz


#
# Stage 2: compile and install giflossy
#
FROM alpine:3.8 AS build-giflossy

RUN apk add --no-cache curl autoconf automake make build-base
RUN curl -SL https://github.com/kornelski/giflossy/archive/1.91.tar.gz | tar xzv
RUN cd giflossy-1.91 && autoreconf -i && ./configure && make install
RUN cp "$(which gifsicle)" /tmp/gifsicle


#
# Stage 2: Compile and build Go app
#

# ...


#
# Stage 3: build final image
#
FROM scratch

COPY --from=builder-ffmpeg /tmp/ffmpeg*/ffmpeg /bin/
COPY --from=builder-giflossy /tmp/gifsicle /bin/

RUN chmod +x /bin/ffmpeg && chmod +x /bin/gifsicle

# ENTRYPOINT /bin/bot-gifv-to-gif
