# Devs Online (DOL)

> [!WARNING]
> This software is in early development, things may be broken or change.

## Overview

This terminal chat app is an attempt to replace Discord, Slack, etc
for developers who want to stay in their terminal and not have a
separate application for teams or friends.

If you are a neovim/vim/emacs user the goal is to have this chat
in your editor.

Currently the chat is just kept in memory so when the server shuts
down the history and room are gone. My next steps are focused
on authentication and choosing a database solution.

The UI is built using the great collection of packages
from [charm](https://charm.sh/)

## Quick Start

Clone the repository and run make build. This will
build the binary as `dol` (Devs Online).

```bash
make build
```

Adding it to user binaries
(Mac Os/Linux)

```bash
mv ./dol /usr/local/bin
```

Windows

```powershell
mv .\dol %USERPROFILE%\bin
```

Running the chat client

```bash
dol chat
```

Running the server

```bash
# local
dol serve

# external
dol serve --external
```

## Road Map

Check out the [GitHub Issues](https://github.com/michael-duren/tui-chat/issues) for current development plans and progress.

## Contributing

If you want to contribute or have an idea for
this application please message me or email
me at `michaeld@michaelduren.com`
