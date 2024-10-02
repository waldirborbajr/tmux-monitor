<img src="https://github.com/user-attachments/assets/9fb07638-5907-4915-b9bf-1ca89255a93d" alt="drawing" style="width:100px;"/>

# Tmux Docker Monitor

<p align="center">
  <img width="256" height="256" src="https://github.com/user-attachments/assets/e9454f13-025f-4a0b-9296-a1807c2cc6c3" />
</p>

A Tmux plugin to monitor Docker containers running on a remote server.

## Features

- Monitor Docker containers on a remote server
- Display the number of containers that are UP, Down, Stopped, Failed, and Died
- Show CPU usage and Memory usage for each container
- Configurable update interval

## Installation

1. Install the Tmux Plugin Manager (TPM) if you haven't already.
2. Add the following line to your `~/.tmux.conf` file:

   ```
   set -g @plugin 'yourusername/tmux-monitor'
   ```

3. Press `prefix + I` to install the plugin.

## Configuration

Create a configuration file at `$HOME/.tmux-monitor` with the following content:

```
SERVER_ADDRESS=your_server_address
SERVER_PORT=22
SERVER_USER=your_username
SERVER_PASSWORD=your_password
UPDATE_INTERVAL=30
```

Replace the placeholder values with your actual server information:

- `SERVER_ADDRESS`: The IP address or hostname of your remote server
- `SERVER_PORT`: The SSH port of your remote server (default is 22)
- `SERVER_USER`: Your SSH username
- `SERVER_PASSWORD`: Your SSH password
- `UPDATE_INTERVAL`: The interval in seconds between updates (default is 30)

## Usage

Once installed and configured, the plugin will automatically start monitoring your Docker containers. The status will be displayed in your tmux status bar.

## Building from Source

To build the plugin from source:

1. Ensure you have Go installed on your system.
2. Clone this repository.
3. Run `make deps` to fetch the required dependencies.
4. Run `make build` to build the plugin.
5. Run `make install` to install the plugin to your tmux plugins directory.

## Development

To contribute to this project:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

- [Tmux](https://github.com/tmux/tmux)
- [Docker](https://www.docker.com/)
- [Go Programming Language](https://golang.org/)
