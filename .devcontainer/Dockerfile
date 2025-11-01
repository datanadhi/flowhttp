# Use an official Go devcontainer base
FROM mcr.microsoft.com/devcontainers/go:1.23

# Install useful tools
RUN apt-get update && apt-get install -y \
    make \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Optional: set up workspace
WORKDIR /workspace/data-nadhi-agent
