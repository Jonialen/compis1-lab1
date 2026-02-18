# Dockerfile for building flex (Fast Lexical Analyzer) from source
# https://github.com/westes/flex

FROM ubuntu:22.04

# Avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install build dependencies
# Note: flex is needed to bootstrap building flex from source
RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    autoconf \
    automake \
    libtool \
    m4 \
    bison \
    flex \
    gettext \
    autopoint \
    texinfo \
    help2man \
    && rm -rf /var/lib/apt/lists/*

# Clone the flex repository
WORKDIR /opt
RUN git clone https://github.com/westes/flex.git

# Build flex from source
# Note: DO NOT use parallel make (-j) - it causes race conditions during bootstrap
WORKDIR /opt/flex
RUN ./autogen.sh \
    && ./configure --prefix=/usr/local \
    && make \
    && make install

# Verify installation
RUN flex --version

# Set up working directory for user files
WORKDIR /workspace

# Keep container running for interactive use
CMD ["sleep", "infinity"]