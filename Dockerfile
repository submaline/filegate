FROM golang:bullseye

RUN apt update && \
    apt install -y git build-essential automake libtool swig libxml2-dev libfftw3-dev \
	             libmagickwand-dev libopenexr-dev liborc-0.4-0 gobject-introspection \
                     libgsf-1-dev libglib2.0-dev liborc-0.4-dev gtk-doc-tools libopenslide-dev \
                     libmatio-dev libgif-dev libwebp-dev libjpeg62-turbo-dev libexpat1-dev

RUN git clone https://github.com/libvips/libvips.git && \
    cd libvips && \
    ./autogen.sh && \
    make && \
    make install

ENV VIPSHOME /usr/local
ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:$VIPSHOME/lib
ENV PATH $PATH:$VIPSHOME/bin
ENV PKG_CONFIG_PATH $PKG_CONFIG_PATH:$VIPSHOME/lib/pkgconfig
ENV MANPATH $MANPATH:$VIPSHOME/man

WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o app
CMD ["./app"]
