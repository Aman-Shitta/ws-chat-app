# WebSocket Chat Server

## Overview
This is a real-time WebSocket-based chat server built using Golang's `golang.org/x/net/websocket` package. The server allows users to join public and private chat channels, send messages, and create new channels dynamically.

## Features
- **WebSocket Support**: Handles real-time chat messages.
- **Multiple Channels**: Users can join or create channels.
- **Broadcast Messaging**: Messages are sent to all users in a channel.
- **Command-Based Interaction**:
  - `SHOW_CHANNELS` - Displays all available channels.
  - `CREATE <channel>` - Creates a new chat channel.
  - `JOIN <channel>` - Joins a specific chat channel.
  - `LEAVE <channel>` - Leaves a specific chat channel.
  - `MESSAGE <channel> <message>` - Sends a message to a channel.
  
## Installation
### Prerequisites
- Go 1.18 or later
- `golang.org/x/net/websocket` package

### Steps
1. Clone the repository:
   ```sh
   git clone https://github.com/Aman-Shitta/ws-chat-app.git
   cd ws-chat-app
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Run the server:
   ```sh
   go run main.go
   ```

## Usage
### Connecting a WebSocket Client
Use a WebSocket client like `wscat` or a simple web-based client to connect:
```sh
wscat -c ws://localhost:8080/ws
```

### Example Commands
```sh
CREATE sports
JOIN sports
MESSAGE sports Hello everyone!
SHOW_CHANNELS
LEAVE sports
```

## Contributing
Feel free to fork this project and submit pull requests to improve the functionality.

## License
This project is open-source and available under the Apache2.0 License.

