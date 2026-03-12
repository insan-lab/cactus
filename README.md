# Cactus

[![Code Triagers Badge](https://www.codetriage.com/furqansoftware/cactus/badges/users.svg)](https://www.codetriage.com/furqansoftware/cactus)

Cactus is a programming contest hosting application. It is a single all-in-one binary that lets you host on-site programming contests over a local area network, with a web-based interface for managing problems, submissions, and standings.

![](screen.png)

## Features

- Single binary with embedded web UI — no client installation needed
- Automatic judging with sandboxed code execution
- Support for C and C++ submissions
- Real-time updates via WebSocket
- Contest standings with scoring and penalty calculation
- Clarification (Q&A) system
- User management with participant, judge, and administrator roles
- CSV account import

## Build

``` sh
make
```

## Usage

``` sh
sudo ./cactus
```

Cactus requires `sudo` for sandboxing submitted code. On first run, it creates a `config.tml` file that you can edit:

``` toml
[core]
addr = ":5050"    # Server address

[belt]
size = 2          # Number of concurrent judging workers
```

The web interface will be available at `http://localhost:5050`.

## Why Cactus?

In 2014, I wanted to build an alternative to PC^2; a better way of hosting on-site programming contests over the local area network. Unlike PC^2, Cactus doesn't need a client to be installed on every computer. Cactus can be distributed as a single all-in-one binary.

Cactus served as my creative outlet until I started working on [Toph](https://toph.co).

## License

Cactus is available under the [BSD 3-Clause License](LICENSE).
