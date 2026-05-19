FROM ubuntu:22.04

# Avoid interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# Use China mirrors for apt
RUN sed -i 's|http://archive.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list && \
    sed -i 's|http://security.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list

# Install dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    git \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Install Go from China mirror
ARG GO_VERSION=1.24.5
RUN curl -fsSL "https://mirrors.aliyun.com/golang/go${GO_VERSION}.linux-amd64.tar.gz" | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"

# Use Go module proxy
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# Install Wails
RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest

WORKDIR /app
