## Задание 1

### Описание доменов, поддоменов и ограниченных конекстов системы (To Be)

- Домен: Управление пользователями
  - Контекст: Аутентификация и авторизация
  - Контекст: Управление пользователями
- Домен: Управление фильмами
  - Контекст: Управление контентом
  - Контекст: Управление метаданными
- Домен: Управление платежами
  - Контекст: Обработка платежей
- Домен: Управление подписками
  - Контекст: Управление подписками

[Container diagram: Кинобездна (To-Be)](diagrams/cinemaabyss_containers_diagram.png)


## Задание 2

### 1. Proxy
- Код находится в каталоге ./src/microservices/proxy.
- Сервис реализован на языке Java при помощи Spring boot.

### 2. Kafka
- Код находится в каталоге ./src/microservices/events.
- Сервис реализован на языке Go.

- [Выполненные тесты](screenshots/ex_2/tests.PNG)
- [Топик: movie-events](screenshots/ex_2/movie-events.PNG)
- [Топик: user-events](screenshots/ex_2/user-events.PNG)
- [Топик: payments-events](screenshots/ex_2/payment-events.PNG)


## Задание 3

### CI/CD

Cборка успешно прошла, все тесты выполнены успешно.<br/> 
В github registry **ghcr.io/zergnk** появились образы:
- architecture-pro-cinemaabyss/proxy-service
- architecture-pro-cinemaabyss/monolith
- architecture-pro-cinemaabyss/movies-service
- architecture-pro-cinemaabyss/events-service

### Proxy в Kubernetes
- [Выполненные тесты](screenshots/ex_3/postman.tests.png)
- [Лог сервиса events-service](screenshots/ex_3/events-service.log.png)
- [Ответ на запрос **https://cinemaabyss.example.com/api/movies**](screenshots/ex_3/https.cinemaabyss.example.com.png)

## Задание 4
- [Развертывания helm](screenshots/ex_4/helm.install.png)
- [Ответ на запрос **https://cinemaabyss.example.com/api/movies**](screenshots/ex_4/https.cinemaabyss.example.com.png)

# Задание 5
Компания планирует активно развиваться и для повышения надежности, безопасности, реализации сетевых паттернов типа Circuit Breaker и канареечного деплоя вам как архитектору необходимо развернуть istio и настроить circuit breaker для monolith и movies сервисов.

```bash

helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update

helm install istio-base istio/base -n istio-system --set defaultRevision=default --create-namespace
helm install istio-ingressgateway istio/gateway -n istio-system
helm install istiod istio/istiod -n istio-system --wait

helm install cinemaabyss .\src\kubernetes\helm --namespace cinemaabyss --create-namespace

kubectl label namespace cinemaabyss istio-injection=enabled --overwrite

kubectl get namespace -L istio-injection

kubectl apply -f .\src\kubernetes\circuit-breaker-config.yaml -n cinemaabyss

```

Тестирование

# fortio
```bash
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.25/samples/httpbin/sample-client/fortio-deploy.yaml -n cinemaabyss
```

# Get the fortio pod name
```bash
FORTIO_POD=$(kubectl get pod -n cinemaabyss | grep fortio | awk '{print $1}')

kubectl exec -n cinemaabyss $FORTIO_POD -c fortio -- fortio load -c 50 -qps 0 -n 500 -loglevel Warning http://movies-service:8081/api/movies
```
Например,

```bash
kubectl exec -n cinemaabyss fortio-deploy-b6757cbbb-7c9qg  -c fortio -- fortio load -c 50 -qps 0 -n 500 -loglevel Warning http://movies-service:8081/api/movies
```

Вывод будет типа такого

```bash
IP addresses distribution:
10.106.113.46:8081: 421
Code 200 : 79 (15.8 %)
Code 500 : 22 (4.4 %)
Code 503 : 399 (79.8 %)
```
Можно еще проверить статистику

```bash
kubectl exec -n cinemaabyss fortio-deploy-b6757cbbb-7c9qg -c istio-proxy -- pilot-agent request GET stats | grep movies-service | grep pending
```

И там смотрим 

```bash
cluster.outbound|8081||movies-service.cinemaabyss.svc.cluster.local;.upstream_rq_pending_total: 311 - столько раз срабатывал circuit breaker
You can see 21 for the upstream_rq_pending_overflow value which means 21 calls so far have been flagged for circuit breaking.
```

Приложите скриншот работы circuit breaker'а

Удаляем все
```bash
istioctl uninstall --purge
kubectl delete namespace istio-system
kubectl delete all --all -n cinemaabyss
kubectl delete namespace cinemaabyss
```
