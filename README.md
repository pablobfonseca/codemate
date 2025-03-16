# Codemate

Codemate is a CLI tool that lets you chat with an AI assistant about your code. It gathers context from your project files and uses AI to provide helpful answers to your coding questions.

## Features

- Chat with an AI assistant about your code
- Automatically collects context from your git repository
- Streams responses in real-time
- Respects .gitignore rules when scanning your codebase

## Installation

```bash
# Clone the repository
git clone https://github.com/pablobfonseca/codemate.git

# Go to the project directory
cd codemate

# Build the project
go build -o codemate main.go

# Or install it globally
go install
```

## Usage

```bash
# Ask a question about your code
codemate chat "How do I implement feature X?"

# Get help
codemate --help
```

## Requirements

- Go 1.24 or higher
- Local [Ollama](https://ollama.ai/) server running with the deepseek-coder model

## Setup Ollama

1. Install Ollama from https://ollama.ai/
2. Pull the deepseek-coder model:
   ```bash
   ollama pull deepseek-coder:6.7b
   ```
3. Start the Ollama server:
   ```bash
   ollama serve
   ```

## License

[MIT](LICENSE)