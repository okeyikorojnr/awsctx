# awsctx

A fast, native Go CLI to switch AWS profiles, inspired by `kubectx`.

`awsctx` reads your local AWS credentials file, lists available profiles, and can launch a sub-shell with `AWS_PROFILE` set to the selected profile.

## Features

- List AWS profiles from `~/.aws/credentials`
- Highlight the current profile when `AWS_PROFILE` is already set
- Switch profiles by starting a sub-shell with the new `AWS_PROFILE`
- Shell completion generation via Cobra

## Requirements

- Go 1.24+
- AWS credentials file at `~/.aws/credentials`

## Install

### Build locally

```bash
go build -o awsctx .
```

Then move it into your PATH, for example:

```bash
sudo mv awsctx /usr/local/bin/
```

### Run without installing

```bash
go run .
```

## Usage

### List profiles

```bash
awsctx
```

This prints all profiles found in `~/.aws/credentials`.
If `AWS_PROFILE` is set, that profile is shown in green.

### Switch profile

```bash
awsctx my-profile
```

If `my-profile` exists, `awsctx` starts a new shell session with:

```bash
AWS_PROFILE=my-profile
```

Exit that shell to return to your previous session.

## Shell completion

Generate completion scripts:

```bash
awsctx completion bash
awsctx completion zsh
```

Current implementation supports `bash` and `zsh` output.

## How it works

- Profiles are read from section names in `~/.aws/credentials`
- `default` and `DEFAULT` are ignored in listing
- A child shell is launched using `$SHELL` (falls back to `/bin/bash`)
- The child shell inherits environment variables and adds `AWS_PROFILE=<target>`

## Notes

- This tool switches profile context only inside the spawned sub-shell.
- It does not modify your AWS config/credentials files.


