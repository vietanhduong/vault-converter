# Vault Converter

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![CI](https://github.com/vietanhduong/vault-converter/actions/workflows/ci.yaml/badge.svg?branch=master&event=push)](https://github.com/vietanhduong/vault-converter/actions/workflows/ci.yaml)

**Support converting Vault Secrets to different formats.**

`vault-converter` is a tool designed to synchronize variables from local to Vault and vice versa.
Currently, `vault-converter` only supports files with the extension `tfvars`.

`vault-converter` uses Vault authentication method as **userpass** with fixed path **userpass/**. But you still can
authorize with **token** method by creat a file contain client token at **"$HOME/.vault_converter/token"**.

Secret Engine supports **Key/Value Version 2** *(kv2)*.

## Installation

### Binaries (recommended)

Download your preferred asset from the [releases page](https://github.com/vietanhduong/vault-converter/releases) and
install manually.

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
$ vault-converter --help
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
$ vault-converter auth --help
Authenticates users to Vault using the provided arguments. 
Method using: 'userpass'. The path of 'userpass' should be 'userpass/'

Usage:
  vault-converter auth [flags]

Flags:
  -a, --address string    Address of the Auth server. This can also be specified via the VAULT_ADDR environment variable. (default "https://dev-vault.knstats.com")
  -h, --help              help for auth
  -p, --password string   The user's password. This can also be specified via the VAULT_PASSWORD environment variables.
  -u, --username string   The username to authenticate with Auth server. This can also be specified via the VAULT_USER environment variables.

Global Flags:
  -v, --version   Print version information and exit. This flag is only available at the global level.
```

### Sync variables from Vault to local

When you pull variables from Vault to local. `vault-convert` automatically override the content to the output file. Keep
it in mind, if you don't want your variables to disappear.

```console
$ vault-converter pull --help
Pull secrets from Vault with specified secret path and convert to file.
SECRET_PATH should be a absolute path at Vault and the values should be in JSON format.
Supports the following formats: "tfvars"

Usage:
  vault-converter pull SECRET_PATH [flags]

Flags:
  -a, --address string   Address of the Auth server. This can also be specified via the VAULT_ADDR environment variable. (default "https://dev-vault.knstats.com")
  -f, --format string    Output format (default "tfvars")
  -h, --help             help for pull
  -o, --output string    Output path. E.g: ~/data/variables.auto.tfvars (default "variables.auto.tfvars")

Global Flags:
  -v, --version   Print version information and exit. This flag is only available at the global level.
```

### Sync variables from local to Vault

Sync variables from local to Vault. If the `SECRET_PATH` doesn't exist. `vault-converter` automatically create new path
and push the content in there. But if the root path *(secret engine path)* does NOT exist, the request will be **fail**
.

```console
$ vault-converter push --help
Parse source file and push secrets to Vault.
Based on the extension of SOURCE_FILE to determine the file format. 
SECRET_PATH should be a absolute path at Vault and the values should 
be in JSON format.
Supports the following formats: "tfvars"

Usage:
  vault-converter push SOURCE_FILE SECRET_PATH [flags]

Flags:
  -a, --address string   Address of the Auth server. This can also be specified via the VAULT_ADDR environment variable. (default "https://dev-vault.knstats.com")
  -h, --help             help for push

Global Flags:
  -v, --version   Print version information and exit. This flag is only available at the global level.
```