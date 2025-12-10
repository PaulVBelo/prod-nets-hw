# Note Service

Микросервис для управления заметками, реализованный на Go.
Поддерживает REST, gRPC, SOAP, масштабируется через несколько экземпляров и работает за HTTPS-балансировщиком (nginx).
Сервис хранит данные в Postgres, также реализовано шардирование.

CRUD-операции:
- Создать заметку
- Получить заметку по ID
- Получить список всех заметок
- Обновить описание заметки
- Удалить заметку

Архитектура:
- Несколько экземпляров сервиса (api-1, api-2)
- За ними — nginx в роли балансировщика нагрузки
- На nginx настроен HTTPS
- Таймауты на обработку: <= 2 секунд
- Корректная обработка ошибок недоступности сервисов и БД
- Шардирование: две postgres db, сущности поровну распределены между ними.


## Public interface

**REST** | base url: https://localhost:8443/notes

| Method   | Endpoint                 | Description       |
| -------- | ------------------------ | ----------------- |
| `POST`   | `/notes`                 | Создать заметку   |
| `GET`    | `/notes`                 | Получить список   |
| `GET`    | `/notes/:id`             | Получить по ID    |
| `PATCH`    | `/notes/:id/description` | Обновить описание |
| `DELETE` | `/notes/:id`             | Удалить           |

Пример:

POST https://localhost:8443/notes
```json
{
    "title": "test title",
    "description": "test description"
}
```

**gRPC**
сервис находится в /api/proto/notes.proto и слушает порт :50051

**SOAP**
Endpoint POST /soap, content-type: text/xml

Пример:
```
<Envelope>
  <Body>
    <CreateNoteRequest>
      <Title>Hello</Title>
      <Description>SOAP example</Description>
    </CreateNoteRequest>
  </Body>
</Envelope>
```
