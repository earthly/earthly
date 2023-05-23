FROM mcr.microsoft.com/devcontainers/base:jammy

# Install Earthly
RUN curl --silent --location https://github.com/earthly/earthly/releases/latest/download/earthly-linux-amd64 --output /usr/local/bin/earthly && \
  chmod +x /usr/local/bin/earthly && \
  /usr/local/bin/earthly bootstrap --with-autocomplete
