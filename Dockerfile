FROM alpine:3.9

#  djmutua/doxmt:latest

LABEL maintainer="Lone Wolf"

ADD $VIMRC_FILE /root
ADD $BASHRC_FILE /root

RUN apk add --update \
    apg \
    bash \
    bind-tools \
    coreutils \
    curl \
    dpkg \
    git \
    g++ \
    gettext \
    jq \
    libc6-compat \
    make \
    nodejs-npm \
    openssh-client \
    openssl \
    openssl-dev \
    python3 \
    py-pip \
    ruby-dev \
    sipcalc \
    tar \
    unzip \
    vim \
    wget
ENV ENV "/root/$BASHRC_FILE"
ENTRYPOINT ["bash","entrypoint.prod.sh"]
