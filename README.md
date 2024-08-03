# Fetch and Run Backend

This is the backend application for the Fetch and Run game, a multiplayer 2D platformer where players collect bones and navigate to the exit while dodging obstacles.

## Requirements

- Go 1.22 or later

## Installation

Clone the repository:
```sh
git clone https://github.com/kevbeltrao/fetch-and-run-backend.git

cd fetch-and-run-backend
```

## Git Hooks

Install Git hooks:
```sh
make prepare
```

Install dependencies:
```sh
make deps
```

## Usage

### Build

Build the application:
```sh
make build
```

### Run

Run the application:
```sh
make run
```

### Test

Run tests:
```sh
make test
```

### Clean

Clean the build files:
```sh
make clean
```

### Format

Format the code:
```sh
make fmt
```

### Lint

Lint the code:
```sh
make lint
```

## Commands

| Command       | Description                         |
|---------------|-------------------------------------|
| `make`        | Build the application               |
| `make run`    | Build and run the application       |
| `make test`   | Run tests                           |
| `make clean`  | Clean the build                     |
| `make fmt`    | Format the code                     |
| `make lint`   | Lint the code                       |
| `make deps`   | Install dependencies                |
| `make help`   | Display help information            |
