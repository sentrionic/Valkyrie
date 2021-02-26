# Valkyrie

<p align="center">
  <img src="https://harmony-cdn.s3.eu-central-1.amazonaws.com/logo.png">
</p>

A [Discord](https://discord.com) clone written in TypeScript.

[Live Demo](https://valkyrieapp.xyz) (Note: File Upload is disabled on the public demo to reduce hosting cost) 

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
- (Basically 2015 Discord features with 2021 Look)

## Stack

- [NestJS](https://nestjs.com/) with [socket.io](https://socket.io/)
- React with [Chakra UI](https://chakra-ui.com/)
- REST Endpoints
- [React Query](https://react-query.tanstack.com/) & [Zustand](https://github.com/pmndrs/zustand) for state management

---

## Installation

### Server

1. Install PostgreSQL and create a DB
2. Install Redis
3. Run `yarn` to install the dependencies
4. Rename `.env.example` to `.env` and fill in the values

- `Required`

        DATABASE_URL="postgresql://<username>:<password>@localhost:5432/db_name"
        REDIS_URL=localhost:6379
        CORS_ORIGIN=http://localhost:3000
        SECRET=SUPERSECRET
        REDIS_HOST=192.168.2.123
        REDIS_PORT=6379
        REDIS_PASSWORD=password

Redis Info is needed twice because the RedisCache Module can't use the `REDIS_URL` directly. 

- `Optional: Not needed to run the app, but you won't be able to upload files or send emails.`

        AWS_ACCESS_KEY=ACCESS_KEY
        AWS_SECRET_ACCESS_KEY=SECRET_ACCESS_KEY
        AWS_STORAGE_BUCKET_NAME=STORAGE_BUCKET_NAME
        AWS_S3_REGION=S3_REGION
        GMAIL_USER=GMAIL_USER
        GMAIL_PASSWORD=GMAIL_PASSWORD

5. Run `yarn start` to run the server

### Web

1. Run `yarn` to install the dependencies
2. Copy .env.example and fill in the values
3. Run `yarn start` to start the client
4. Go to `localhost:3000`

## Endpoints

Once the server is running go to `localhost:4000/api` to see all the HTTP endpoints
and `localhost:4000/ws` for all the websocket events.
