FROM golang:1.20

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=${USER_UID}

RUN groupadd --gid ${USER_GID} ${USERNAME} && \
    useradd --uid ${USER_UID} --gid ${USER_GID} -m ${USERNAME} && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    git \
    groff \
    less \
    sudo \
    unzip && \
    echo ${USERNAME} ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/${USERNAME} && \
    chmod 0440 /etc/sudoers.d/${USERNAME} && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

USER ${USERNAME}

# Dev tools
RUN go install github.com/ramya-rao-a/go-outline@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install github.com/kisielk/errcheck@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /workspace
