# Social Network
## Development Workflow with Makefile

The Makefile helps you start/stop services, seed the dev database, and clear/reset databases for local development. It works on Windows (Git Bash/WSL), Linux, and Mac.

## Prerequisites

- Docker & Docker Compose installed

- Services (users, forum, chat, notifications) and databases defined in docker-compose.yml

- services/users/seeds/seed.sh exists for seeding dev data

- DATABASE_URL set correctly if running seed from host

## Makefile Overview

Available targets:

|Target|	Description|
|--|--|
|make build	|Build all images|
|make rebuild	|Run down + build + up|
|make up	|Start all containers in detached mode|
|make down|	Stop all containers|
|make seed|	Seed the dev database (currently seeds users only)|
|make clear|	Clear all databases and restart containers to rerun migrations|
|make clear-users|	Clear users DB and restart users container (migrations run automatically)|
|make clear-forum|	Clear forum DB and restart forum container|
|make clear-chat|	Clear chat DB and restart chat container|
|make clear-notifications	|Clear notifications DB and restart notifications container|
|make reset-users	|Clear users DB and seed it in one command (convenience target)|

## Typical Dev Workflow

### Build all images 
```make build```

### Start all services 
```make up```


### Seed the users database

```make seed```


### Clear all databases (optional)

```make clear```

After make clear, the containers are restarted and migrations are applied automatically.

You can then seed the database again:

```make seed```


### Stop all services

```make down```

### Reset Users DB in One Command

If you want to clear the users database and seed it immediately, use the reset-users target:

```make reset-users```


This runs the following steps automatically:

Clear the users database (DROP SCHEMA public CASCADE; CREATE SCHEMA public)

Restart the users container so migrations are applied

Seed the dev data

### Stop containers, rebuild images and run containers

```make rebuild```

## Notes / Best Practices

Do not run seeds in production.

The Makefile uses docker exec -T for cross-platform compatibility.

sleep 5 is included after container restarts to allow migrations to finish — increase if your migrations are long.

Seeds are idempotent (ON CONFLICT DO NOTHING) so they can be safely rerun.

For other services (forum, chat, notifications), similar seeding scripts will be added as needed.





## Access databases with
```
docker-compose exec users-db psql -U postgres -d social_users
```

```
docker-compose exec chat-db psql -U postgres -d social_chat
```

```
docker-compose exec forum-db psql -U postgres -d social_forum
```

```
docker-compose exec forum-db psql -U postgres -d social_notifications
```

