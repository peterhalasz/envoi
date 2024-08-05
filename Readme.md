# envoi - cloud workstation manager

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

    This will generate an executable file in the current directory.

3. Run the executable:

    ```shell
    ./envoi
    ```

## Usage

### Commands

* `init` - Initializes your workstation
    * Creates a new droplet
    * Sets up a new SSH key if none exists
    * Creates and attaches a volume (if configured) to the droplet
* `delete` - Deletes your workstation
    * Deletes the droplet
    * Unattaches the volume and deletes it (if configured)
* `stop` - Stop your workstation
    * Deletes the droplet only
* `start` - Start your workstation
    * Creates a new droplet
    * Attaches the existing volume (if configured) to the droplet
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

This project is licensed under the [MIT License](LICENSE).
