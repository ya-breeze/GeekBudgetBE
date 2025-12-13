# Multi-language development environment for Angular (TypeScript/Node.js) and Golang
# Based on Microsoft's official devcontainer base images
FROM mcr.microsoft.com/devcontainers/go:dev-1

ENV USERNAME=vscode
ENV NODE_VERSION=22
ENV GOPATH=/go

# Install system dependencies for both Node.js and Go development
RUN apt-get update && apt-get install -y \
    bash-completion \
    sqlite3 \
    curl \
    wget \
    git \
    build-essential \
    python3 \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js and npm for Angular development
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest && \
    npm install -g @angular/cli@latest && \
    npm install -g typescript@latest && \
    rm -rf /var/lib/apt/lists/*

# Install Go development tools
RUN go install golang.org/x/tools/cmd/goimports@latest && \
    chown -R $USERNAME /go

# Set up proper PATH for both Node.js and Go
ENV PATH="/home/${USERNAME}/.local/bin:${GOPATH}/bin:${PATH}"

# Bash history
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/command_history/.bash_history" \
    && mkdir /command_history \
    && touch /command_history/.bash_history \
    && chown -R $USERNAME /command_history \
    && echo "$SNIPPET" >> "/home/$USERNAME/.bashrc"
