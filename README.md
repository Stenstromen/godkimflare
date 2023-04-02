# GoDKIMFlare

Go binary for creating/updating DKIM records on Cloudflare.

<br>

## Generate Cloudflare API Token
1. Visit [https://dash.cloudflare.com/profile/api-tokens](https://dash.cloudflare.com/profile/api-tokens)
2. Create Token
3. "Edit Zone DNS" Template
4. "Zone Resources" Include > Specific Zone > example.com

## Installation via Homebrew (MacOS/Linux - x86_64/arm64)
```
brew install stenstromen/tap/godkimflare
```
## Download and Run Binary
* For **MacOS** and **Linux**: Checkout and download the latest binary from [Releases page](https://github.com/Stenstromen/godkimflare/releases/latest/)
* For **Windows**: Build the binary yourself.

## Build and Run Binary
```
go build
./godkimflare
```

## Example Usage
```
- Create MTA-STS DNS Records and Nginx Configuration
export TOKEN="# Cloudflare API TOKEN"
./godkimflare create -d example.com -f /path/to/privateKeyFile

- Update MTA-STS DNS Records and/or Nginx Configuration
export TOKEN="# Cloudflare API TOKEN"
./godkimflare update -d example.com -f /path/to/privateKeyFile

Go binary for creating/updating DKIM records on Cloudflare.

Usage:
  godkimflare [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create DKIM Record on Cloudflare
  help        Help about any command
  update      update DKIM Record on Cloudflare

Flags:
  -h, --help   help for godkimflare

Use "godkimflare [command] --help" for more information about a command.
```

<br>

# Random notes

```
openssl rsa -in dkim_private.pem -pubout -outform der 2>/dev/null | openssl base64 -A
```