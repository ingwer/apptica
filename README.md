###### Тестовый проект

Собрать

`make build`

Запустить

`make up`

Остановить 

`make down`

Потыкать http

`curl '127.0.0.1/appTopCategory?date=2021-06-04'`

Потыкать grpc 

`docker run --network="host" --rm -it networld/grpcurl /grpcurl -plaintext -d '{"date":"2021-06-04"}' 127.0.0.1:9000 AppticaService.AppTopCategories`