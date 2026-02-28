#!/usr/bin/env -S just --justfile
# ^ A shebang isn't required, but allows a justfile to be executed
#   like a script, with `./justfile test`, for example.

# Default task: list all recipes
default:
    just --list

# Install task: set up the repository
install:
    hk install