msvc-demo

Структура каталогов
[common] - общие файлы для разных модулей
[entity-db-postgres] - сборка контейнера alisin69/msvc-pgdb
[entity-service] - модуль заданного сервиса. Там же сборка контейнера alisin69/msvc-entity-service
[entity-web] - модуль web UI для заданного сервиса. Там же сборка контейнера alisin69/msvc-entity-web
[vendor] - все сторонние пакеты выкаченные утилитой DEP

docker-compose.yaml - файл для управления контейнерами.

# Файлы с манифестами для kubernetes
[k8s]/[helm] - генерация файлов для Kubernetes
[k8s]/[orig] - сгенерированный helm файл для Kubernetes

Makefile & Dockerfile контейнеров находятся в 
  [entity-db-postgres]
  [entity-service]
  [entity-web]
  
URL сервисов:
[GET]    /-/healthy      Check liveness
[GET]    /-/ready        Check readiness
[GET]    /-/metrics      Metrics

[GET]    /entities       Get all data
[GET]    /entities/{id}  Get entity with required ID
[POST]   /entities       Create entity
[PUT]    /entities/{id}  Modify entity with required ID
[DELETE] /entities       Delete entity with required ID

Тесты:
\entity-service\cmd\handlers_test.go       - mock  тесты на http handler-ы 
\entity-service\pkg\repository\rep_test.go - модульный тест репозитория при помощи имитации БД (сохранение в память)


To work:
curl -d '{"name":"name 1", "descr":"descr for name 1", "last-operator":"User 1"}' -H "Content-Type: application/json" -X POST http://localhost:9090/entities
curl -d '{"name":"name 1 changes are not affected", "descr":"descr for name 1 changed", "last-operator":"User 2"}' -H "Content-Type: application/json" -X PUT http://localhost:9090/entities/1
curl -H "Content-Type: application/json" -X DELETE http://localhost:9090/entities/1
curl -H "Content-Type: application/json" -X GET http://localhost:9090/entities/1
curl -H "Content-Type: application/json" -X GET http://localhost:9090/entities


=============================
Примеры команд, которые могут пригодиться

dep ensure -v
dep ensure -update  
dep status


docker build -t pgdb .
docker rm $(docker ps -aq)
docker rmi msvc-pgdb msvc-entity-service

go run !(*_test).go

cd entity-service/
go test -v


helm create entity

helm install --dry-run --debug -n xxx ./k8s/helm/entity
helm lint --strict ./k8s/helm/entity
helm delete --dry-run --debug xxx 

helm install -n xxx ./k8s/helm/entity   
helm upgrade --set rest.replicaCount=3 xxx ./k8s/helm/entity 
helm delete --purge xxx 

helm ls --all   

