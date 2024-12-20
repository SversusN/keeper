
![Frame 1 (3)](https://github.com/user-attachments/assets/2ad8c533-4cce-4f9c-8f32-0ce92894df08)

# Keeper
Сервис для хранения и обработки паролей

# Server

Представляет собой gRPC сервер на GO


## Функции сервера

- Аутентификация и авторизация пользователя по логину и паролю. (JWT)
- Сохранение приватных данных, таких как:
    - Пароли;
    - Банковские карты;
    - Файлы;
    - Текстовые данные.
- Просмотр информации о файлах пользователя;
- Получение и просмотр данных;
- Редактирование сохранённых данных;
Все данные представляются как абстрактная "запись", хранятся в бинарном представлении.


## Запуск сервера
Запуск возможен в докере. Созданы doker-файлы для сборки, запуска и работы базы данных Postgres.

# Client

Консольное приложение, выполненное в стиле терминалов 80-ых.
Пользователю предлается сценарий действий, для выполнения функций.

1. Пользователь регистрируется или авторизуется при наличии;
2. После авторизации доступно:
   - Простмотр информации о данных;
   - Запрос детали записи (дешифрованнное);
   - Измнение имеющейся записи;
   - Внесение новой записи.
3.Каждое действие определяется своим сценарием работы.
## Функции клиента
Клиент выдает по требованию данные в текстовом представлении
Клиент может работать оффлайн, при условии что была проведена синхронизация и сохранение в in-memory кэш.
Клиент шифрует и дешифрует данные. Открытая информация не хранится на сервере (кроме общедоступной информации)
Возможна работа нескольких клиентов.

Примерная схема взаимодействия

![keeper-scheme](https://github.com/user-attachments/assets/90a30c4e-edb7-494e-b0e0-5049a3e72f7c)

## Демонстрация работы клиента.


![keeper](https://github.com/user-attachments/assets/fb580e24-5605-4728-9708-28571f29091f)


## Сборка клиента

    # build linux
    go build -ldflags "-X main.buildDate=${current_time} -X main.buildVersion=${git_hash}" -o ./cmd/client/build/linux ./cmd/client/
    # build windows
    GOOS=windows GOARCH=amd64 go build -ldflags"-X main.buildDate=${current_time} -X main.buildVersion=${git_hash}" -o ./cmd/client/build/win ./cmd/client
    # build mac amd64
    GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.buildDate=${current_time} -X main.buildVersion=${git_hash}" -o ./cmd/client/build/mac ./cmd/client
 
