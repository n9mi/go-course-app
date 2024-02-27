# go-course-app
Simulating online learning platform requirements that provided in test.pdf (80% completed)

## **Packages used**
- github.com/spf13/viper
- github.com/gofiber/fiber/v2
- github.com/go-playground/validator/v10
- gorm.io/driver/postgres
- gorm.io/gorm 
- github.com/stretchr/testify
- github.com/redis/go-redis/v9
- github.com/golang-jwt/jwt/v5
- github.com/sirupsen/logrus

## **How to run the app**
```
docker compose build
docker compose up
go run cmd/web/main.go
```

## **Run the migration**
Migration will automatically run when the server 

## **Structure**
Based on repository pattern, this project use:
- Repository layer: For accessing db in the behalf of project to store/update/delete data
- Usecase layer: Contains set of logic/action needed to process data/orchestrate those data
- Entity: Contains set of database atribute
- Model: Contains set of data that will be parsed or send as request or response
- Controller layer: Acts to mapping users input/request and presented it back to user as relevant responses

## TODO
- Refresh endpoint
- Logout endpoint
- Statistic for admin


