# luckysix

## Start docker
```bash
docker-compose up -d
# if u need to rebuild all docker containers
docker-compose up -d --force-recreate
```
Output u can see in Docker Dashboard

## Init, update and install
```bash
docker-compose up -d
go mod init github.com/basel-ax/luckysix
go get -u ./... && go mod tidy 
# just install
go mod tidy
```

## Start project 
Localy, with enviroment in docker containers.  
```bash
go run ./cmd/app
```

## Project Info
Application Based on [Go-clean-template](https://github.com/evrone/go-clean-template)  
Sheduller Based on [jasonlvhit/gocron](https://github.com/jasonlvhit/gocron)  
Add to project - go get github.com/jasonlvhit/gocron  

Telegram API [Go Telegram Bot API](https://go-telegram-bot-api.dev/)  
Add to project - go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5 

## ORM
[GORM](https://gorm.io/), top [from list](https://github.com/d-tsuji/awesome-go-orms) 
