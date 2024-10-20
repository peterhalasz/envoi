# envoi - cloud workstation manager

## Warning
**This project is currently under development and should not be relied on for any kind of production or critical processes.** There is no guarantee that the resources will be created or deleted as expected. Please use with caution.

## Description

Sometimes, I need a small Linux VM to test long-running or potentially nasty software.
However, I kept running into issues on my M1 Mac with virtualization. I'm a fan of DigitalOcean but I don't like clicking around on their GUI. The DO cli is great, but I wanted something smaller and simpler.

Behold, envoi.
With envoi you can easily manage your cloud workstation on DigitalOcean. The application supports 
- creating
- deleting
- connecting to
virtual machines on Digital Ocean.

## Building

To build this project, you can follow these steps:

1. Make sure you have Go installed on your machine. You can download it from the official Go website: https://golang.org/dl/

2. Build the project using the `go build` command:

    ```shell
    go build
    ```

    This will generate the executable in the current directory.

3. Run the executable:

    ```shell
    ./envoi
    ```

## Usage

### Setup
You need to have a valid DigitalOcean token to use envoi. You can generate one under [API->Tokens](https://cloud.digitalocean.com/account/api/tokens)
Once you have the token you need to set up an environment variable with `export DO_TOKEN=<YOUR API TOKEN>`.
Alternatively paste the token into a file and set up a configuration parameter for `digitalocean.token_path`.

Envoi is going to use your SSH public key from `~/.ssh/id_rsa.pub`. This is configurable with the param `ssh.public_key_path`.

### Commands

* `start` - Starts your workstation
    * Creates a new droplet
    * Sets up a new SSH key if none exists
    * Creates and attaches a volume (if configured) to the droplet
* `delete` - Deletes your workstation
    * Deletes the droplet
    * Unattaches the volume and deletes it (if configured)
* `stop` - Stop your workstation
    * Deletes the droplet only
    * Keeps the volume (if configured)
* `status` - Prints the status of your workstation
* `connect` - Connects to the workstation via ssh

### Config
You can configure a few things by creating a new file `.envoi.conf` and placing it next to the `envoi` executable or in your home directory.
You can find all possible configuration keys under `internal/util/config.go`.

## Contributing

Contributions are welcome!
Feel free to open a PR to fix smaller issues.
If you have something bigger in mind, please open an issue and we can discuss it.

## License

This project is licensed under the [MIT License](LICENSE.txt).
