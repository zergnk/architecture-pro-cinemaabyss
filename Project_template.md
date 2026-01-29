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
Реализация сервиса proxy-service
- Код находится в каталоге ./src/microservices/proxy.
- Сервис реализован на языке Java при помощи Spring boot.

### 2. Kafka
Реализация сервиса events-service
- Код находится в каталоге ./src/microservices/events.
- Сервис реализован на языке Go.

Скриншоты логов и выполненных тестов:
- [Выполненные тесты](screenshots/ex_2/tests.PNG)
- [Топик: movie-events](screenshots/ex_2/movie-events.PNG)
- [Топик: user-events](screenshots/ex_2/user-events.PNG)
- [Топик: payments-events](screenshots/ex_2/payment-events.PNG)


## Задание 3

### CI/CD
- **api-tests.yml** выполняются автоматически после push в ветках main и cinema. Также он может быть запущен в ручном режиме
в обеих ветках
- **docker-build-push.yml** выполняется автоматически после push только в ветке main. Но может быть запущен в ручном режиме 
- в ветках main и cinema

Cборка прошла без ошибок, все тесты выполнены успешно.<br/> 
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
- [Лог fortio](screenshots/ex_5/fortio-log.png)
