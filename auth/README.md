# auth

## Usage

- Start app with docker-compose (with logs option): \
```docker-compose --env-file .\build\docker\configs\prod.env up --build > compose-logs.log```

- Start just a postgres with docker-compose and other services in determined shell (with logs option):
```docker-compose --env-file .\build\docker\configs\prod.env up --build > compose-logs.log```