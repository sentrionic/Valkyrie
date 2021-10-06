[![Go Report Card](https://goreportcard.com/badge/github.com/sentrionic/Valkyrie)](https://goreportcard.com/report/github.com/sentrionic/Valkyrie)
[![Netlify Status](https://api.netlify.com/api/v1/badges/cd1667ed-3257-41d0-82ca-7b34de655339/deploy-status)](https://app.netlify.com/sites/valkyrie-app/deploys)

# Valkyrie

<p align="center">
  <img src="https://harmony-cdn.s3.eu-central-1.amazonaws.com/logo.png">
</p>

A [Discord](https://discord.com) clone using [React](https://reactjs.org/) and [Go](https://golang.org/).

[Live Demo](https://valkyrieapp.xyz)

**Notes:**

- File Upload is disabled.
- The design does not fully match current Discord anymore.
- Data regularly gets wiped, so you can use any valid email and password.
- For the old [Socket.io](https://socket.io/) stack using [NestJS](https://nestjs.com/) check out the [v1](https://github.com/sentrionic/Valkyrie/tree/v1) branch.

## Video

![Preview](.github/preview.gif)

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

### Server

- [Gin](https://gin-gonic.com/) for the HTTP server
- [Gorilla Websockets](https://github.com/gorilla/websocket) for WS communication
- [Gorm](https://gorm.io/) as the database ORM
- PostgreSQL to save all data
- Redis for storing sessions and reset token
- S3 for storing files and Gmail for sending emails
- Hosted on [Heroku](https://www.heroku.com/)

### Web

- React with [Chakra UI](https://chakra-ui.com/)
- [React Query](https://react-query.tanstack.com/) & [Zustand](https://github.com/pmndrs/zustand) for state management
- [Typescript](https://www.typescriptlang.org/)
- Hosted on [Netlify](https://www.netlify.com/)

For the mobile app using Flutter check out [ValkyrieApp](https://github.com/sentrionic/ValkyrieApp)

---

## Installation

### Server

If you are familiar with `make`, take a look at the `Makefile` to quickly setup the following steps
or alternatively copy the commands into your CLI.

1. Install Docker and get the Postgresql and Redis containers (`make postgres` && `make redis`)
2. Start both containers (`make start`) and create a DB (`make createdb`)
3. Install Golang and get all the dependencies (`go mod tidy`)
4. Rename `.env.example` to `.env` and fill in the values

- `Required`

        PORT=4000
        DATABASE_URL=postgresql://<username>:<password>@localhost:5432/valkyrie
        REDIS_URL=redis://localhost:6379
        CORS_ORIGIN=http://localhost:3000
        SECRET=SUPERSECRET
        HANDLER_TIMEOUT=5
        MAX_BODY_BYTES=4194304 # 4MB in Bytes = 4 * 1024 * 1024

- `Optional: Not needed to run the app, but you won't be able to upload files or send emails.`

        AWS_ACCESS_KEY=ACCESS_KEY
        AWS_SECRET_ACCESS_KEY=SECRET_ACCESS_KEY
        AWS_STORAGE_BUCKET_NAME=STORAGE_BUCKET_NAME
        AWS_S3_REGION=S3_REGION
        GMAIL_USER=GMAIL_USER
        GMAIL_PASSWORD=GMAIL_PASSWORD

5. Run `go run github.com/sentrionic/valkyrie` to run the server

**Alternatively**: If you only want to run the backend without installing Golang and all dependencies, you can download the pre compiled server from the [Release tab](https://github.com/sentrionic/Valkyrie/releases) instead. You will still need to follow the above steps 1, 2 and 4.

### Web

0. Install the latest or LTS version of Node.
1. Run `yarn` to install the dependencies
2. Rename `.env.development.example` to `.env.development`.
3. Run `yarn start` to start the client
4. Go to `localhost:3000`

## Endpoints

Once the server is running go to `localhost:<PORT>/swagger/index.html` to see all the HTTP endpoints
and `localhost:<PORT>` for all the websockets events.

## Tests

### Server

All routes in `handler` have tests written for them.

Function calls in the `service` directory that do not just delegate work to the repository have tests written for them.

Run `go test -v -cover ./service/... ./handler/...` (`make test`) to run all tests

Additionally this repository includes E2E tests for all successful requests. To run them you
have to have Postgres and Redis running in Docker and then run `go test github.com/sentrionic/valkyrie` (`make e2e`).

## Credits

[Ben Awad](https://github.com/benawad): The inital project is based on his Slack tutorial series and I always look at his repositories for inspiration.

[Jacob Goodwin](https://github.com/JacobSNGoodwin/memrizr): This backend is built upon his tutorial series and uses his backend structure

[Jeroen de Kok](https://dev.to/jeroendk/building-a-simple-chat-application-with-websockets-in-go-and-vue-js-gao): The websockets structure is based on his tutorial
