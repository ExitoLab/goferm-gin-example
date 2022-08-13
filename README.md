### Goferm demonistration of Go Gin Web Framework

This is a simple demonstration on how use Gin Webframe work 

The application is created using: 

1. Go - Gin Framework 
2. Mongodb - database

You can run the app by following the command below:

## To deploy 

1. Install MongoDB
2. Install golang 
3. Run command `go build -o new -v`
4. Run command `./new` 
5. To run in background use this command `nohup ./new &`

The application is running on port 8001


## Kill running app on port 8001

To kill the app while debugging run this command `kill -9 $(lsof -t -i:8001)`