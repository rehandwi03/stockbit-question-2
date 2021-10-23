## 🏁 Getting Started <a name = "getting_started"></a>

### Prerequisites

First you must install docker & docker-compose on your PC/Laptop.

### Run in Docker Compose

You can run apps with docker compose and make sure you are now in the project folder of apps

```
cp .env.example .env
docker-compose up -d
```

It will run mariadb database and apps in docker container. Note: please adjust value in .env with your environment

## 🔧 Running the tests <a name = "tests"></a>

### Integration Testing

It will test API in app using real database.

```
cd tests
go test ./...
```

## ⛏️ Built Using <a name = "built_using"></a>

- [Mariadb](https://www.mariadb.org/) - Database
- [Golang](https://golang.org/) - Server Apps

## ✍️ Authors <a name = "authors"></a>

- [@rehandwi03](https://github.com/rehandwi03) - Idea & Initial work


