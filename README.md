# Vault Converter

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![CI](https://github.com/vietanhduong/vault-converter/actions/workflows/ci.yaml/badge.svg?branch=master&event=push)](https://github.com/vietanhduong/vault-converter/actions/workflows/ci.yaml)

**Support converting Vault Secrets to different formats.**

`vault-converter` is a tool designed to synchronize variables from local to Vault and vice versa.
Currently, `vault-converter` only supports files with the extension `tfvars`.

`vault-converter` uses Vault authentication method as **userpass** with fixed path **userpass/**. 
But you still can authorize with **token** method by creat a file contain client token at **"$HOME/.vault_converter/token"**.

Secret Engine supports **Key/Value Version 2** *(kv2)*. 

## Installation

### Binaries (recommended)

Download your preferred asset from the [releases page](https://github.com/vietanhduong/vault-converter/releases) and install manually.

### Source code

```shell
# clone repo to some directory outside GOPATH
git clone https://github.com/vietanhduong/vault-converter

cd vault-converter

go mod download

go build . 
```

## Usage 

Currently, `vault-converter` supports synchronize variables from Vault to local and vice versa.

```console
Convert to file from Vault. Support multiple file format like '.tfvars', '.env'

Usage:
  vault-converter [flags]
  vault-converter [command]

Available Commands:
  auth        Authenticates users to Vault
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  pull        Pull secrets from Vault and convert to file
  push        Parse source file and push to Vault

Flags:
  -h, --help      help for vault-converter
  -v, --version   Print version information and exit. This flag is only available at the global level.

Use "vault-converter [command] --help" for more information about a command.
```

### Authenticate
User authentication with Vault 
```console

```