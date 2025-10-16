FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive \
    TZ=Etc/UTC \
    NVM_DIR=/root/.nvm

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    wget \
    git \
    gnupg \
    software-properties-common \
    build-essential \
    make \
    gcc \
    g++ \
    clang \
    cmake \
    libmysqlclient-dev \
    libssl-dev \
    libbz2-dev \
    libreadline-dev \
    libncurses-dev \
    libboost-all-dev \
    p7zip-full \
    pkg-config \
    sudo \
    openssh-client \
    vim \
    less \
    && rm -rf /var/lib/apt/lists/*

RUN update-alternatives --install /usr/bin/cc cc /usr/bin/clang 100 && \
    update-alternatives --install /usr/bin/c++ c++ /usr/bin/clang 100

SHELL ["/bin/bash", "-lc"]

RUN mkdir -p "$NVM_DIR" \
    && git clone https://github.com/nvm-sh/nvm.git "$NVM_DIR" \
    && cd "$NVM_DIR" && git checkout v0.34.0 \
    && source "$NVM_DIR/nvm.sh" \
    && nvm install 20.18.0 \
    && nvm alias default 20.18.0

ENV PATH="$NVM_DIR/versions/node/v20.18.0/bin:$PATH"

RUN mkdir -p /tswow-root/source

COPY docker/tswow-entrypoint.sh /usr/local/bin/tswow-entrypoint.sh
COPY docker/tswow-build.sh /docker/tswow-build.sh
RUN chmod +x /usr/local/bin/tswow-entrypoint.sh /docker/tswow-build.sh

WORKDIR /tswow-root

ENTRYPOINT ["/usr/local/bin/tswow-entrypoint.sh"]
CMD ["bash"]
