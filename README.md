# symlinker

**symlinker** is a command-line tool for managing symbolic links (symlinks) on your filesystem using a simple YAML configuration file. It allows you to create, verify, and delete symlinks in a safe and repeatable way.

## Features

- **Declarative configuration:** Define all your symlinks in a single YAML file.
- **Safe operations:** Never overwrites or deletes files that are not symlinks.
- **Cross-platform:** Works anywhere Go and symlinks are supported.
- **Idempotent:** Running the same command multiple times is safe.

## Installation

Install via Homebrew (recommended for macOS / Linux with Homebrew):

```sh
# add the tap, then install
brew tap thecheerfuldev/cli
brew install thecheerfuldev/cli/symlinker
```

Manual download:

- Visit the releases page and download the appropriate archive for your platform:
  https://github.com/TheCheerfulDev/symlinker/releases/latest
- Downloads available for:
  - macOS: amd64, arm64
  - Linux: amd64, arm64
- Extract the archive and place the `symlinker` binary on your PATH. For example:
  ```sh
  tar -xzf symlinker_<version>_<os>_<arch>.tar.gz
  mv symlinker /usr/local/bin/symlinker
  chmod +x /usr/local/bin/symlinker
  ```

## Usage

### 1. Initialize a configuration file

Create a default `symlinker.yaml` in your current directory:

```sh
symlinker init
```

This creates a file like:

```yaml
symlinks:
  - source: /path/to/source
    target: /path/to/target
```

Edit this file to list all the symlinks you want to manage.

### 2. Apply symlinks

Create all symlinks defined in your configuration:

```sh
symlinker apply
```

- Skips links if the source does not exist.
- Verifies existing symlinks.
- Never overwrites or deletes existing files.

### 3. Verify symlinks

Check that all symlinks are present and correct:

```sh
symlinker verify
```

- Reports missing sources, targets, or mismatched symlinks.

### 4. Delete symlinks

Remove symlinks defined in your configuration:

```sh
symlinker delete
```

- Only deletes symlinks that point to the specified source.
- Does not delete regular files or directories.

### 5. Custom configuration file

By default, `symlinker` uses `symlinker.yaml` in the current directory. To use a different file:

```sh
symlinker --file mylinks.yaml apply
```

## Configuration File Format

The YAML file should look like:

```yaml
symlinks:
  - source: /absolute/or/relative/path/to/source
    target: /absolute/or/relative/path/to/target
  # Add more links as needed
```

- `~` is expanded to your home directory.
- Relative paths are resolved from the current working directory.

## License

Apache License 2.0

## Contributing

Pull requests and issues are welcome!
