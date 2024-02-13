# Card game server

This repository contains the code for the card game server in Golang. The database used is redis contolled using go redis. The user can register/login with a secret. The user can start the game .

## Features

-   Save game status
-   Get leaderboard
-   Verify user login
-   Redister user
-   Get game data of user

## Run Locally

Clone the project

```bash
  git clone https://github.com/Balaguru1601/card-game-server.git
```

Go to the project directory

```bash
  cd card-game-server
```

Install dependencies

```bash
  go get
```

```bash
  go install github.com/githubnemo/CompileDaemon
```

Start the server

```bash
  CompileDaemon -command="./go-backend"
```

The server will start at port 8080.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`REDIS_URL`

`PORT`
