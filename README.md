# Concrnt CLI
A command line interface client for the [concrnt](https://github.com/concrnt) social media platform. Written in Go. This is a very WIP project.

<img src="./concrnt-cli.gif" width="600">

# Running
Go is needed to run this program. You can download it [here](https://golang.org/dl/) or install it via your package manager.

## Configuration
The configuration is stored in `.env` file. Filling out the `.env` file is required to run the program. The following are the required fields:
- `CCID`: Your Concrnt CCID (public key)
- `PRIVATE_KEY`: Your Concrnt master private key or subkey
- `CKID`: Your Concrnt CKID (in case of using a subkey)
- `TIMELINES`: The timelines to posts to. Separated by commas. The fetch function will only fetch from the first timeline for now.
