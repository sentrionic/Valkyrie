# Valkyrie

<p align="center">
  <img src="https://harmony-cdn.s3.eu-central-1.amazonaws.com/logo.png">
</p>

A [Discord](https://discord.com) clone using React and Go.

[Live Demo](https://valkyrieapp.xyz)

**Notes:**

- File Upload is disabled.
- The live demo currently runs on the [Go backend](https://github.com/sentrionic/ValkyrieGo).
- Design does not fully match current Discord anymore.
- Data regularly gets wiped, so you can use any valid email and password.

## Video

![Game](preview.gif)

## Features

- Message, Channel, Server CRUD
- Authentication using Express Sessions
- Channel / Websocket Member Protection
- Realtime Events
- File Upload (Avatar, Icon, Messages) to S3
- Direct Messaging
- Private Channels
- Friend System
- Notification System
- Basic Moderation for the guild owner (delete messages, kick & ban members)

## Stack

- [Gin](https://gin-gonic.com/) for the HTTP server
- [Gorilla Websockets](https://github.com/gorilla/websocket) for WS communication
- React with [Chakra UI](https://chakra-ui.com/)
- REST Endpoints
- [React Query](https://react-query.tanstack.com/) & [Zustand](https://github.com/pmndrs/zustand) for state management

For the mobile app using Flutter check out [ValkyrieApp](https://github.com/sentrionic/ValkyrieApp)

---

## Installation

### Server

Go to the [ValkyrieGo](https://github.com/sentrionic/ValkyrieGo) repository and follow the instructions there.

### Web

0. Install the latest or LTS version of Node.
1. Run `yarn` to install the dependencies
2. Rename `.env.development.example` to `.env.development`.
3. Run `yarn start` to start the client
4. Go to `localhost:3000`

## Credits

[Ben Awad](https://github.com/benawad): The inital project is based on his Slack tutorial series and I always look at his repositories for inspiration.
