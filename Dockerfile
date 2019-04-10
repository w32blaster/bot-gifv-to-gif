#
# Stage 1: Install dependencies
#
FROM alpine AS builder

ARG FFMPEG_VERSION=4.1.3
ARG FFMPEG_URL=https://johnvansickle.com/ffmpeg/releases/ffmpeg-$FFMPEG_VERSION-64bit-static.tar.xz

# Download ffmpeg
ADD ${FFMPEG_URL} /tmp/ffmpeg.tar.xz
RUN cd /tmp && tar xJf ffmpeg.tar.xz

#
# Stage 2: Compile and build Go app
#

# ...


#
# Stage 3: build final image
#
FROM scratch

COPY --from=builder /tmp/ffmpeg*/ffmpeg /bin/
RUN chmod +x /bin/ffmpeg

# ENTRYPOINT /bin/bot-gifv-to-gif
