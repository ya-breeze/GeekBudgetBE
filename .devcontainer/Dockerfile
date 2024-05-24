FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

RUN apt-get update && apt-get install \
    bash-completion

ENV USERNAME=vscode

# Bash history
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/command_history/.bash_history" \
    && mkdir /command_history \
    && touch /command_history/.bash_history \
    && chown -R $USERNAME /command_history \
    && echo "$SNIPPET" >> "/home/$USERNAME/.bashrc"
