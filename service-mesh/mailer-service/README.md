# Mailer Service

Отвечает за отправку заметок по электронной почте.
Получает данные о заметке из note-service и отправляет письмо через Mailhog.

## API

### **POST /send**

отправка письма на указанный email

POST https://localhost:8443/mailer/send

### Request body
```json
{
  "note_id": "some UUID",
  "email": "user@example.com"
}
```

### Response
```json
{
  "status": "sent",
  "to": "user@example.com",
  "note": {
    "id": "...",
    "title": "...",
    "description": "...",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

Посмотреть отправленные письма можно в UI Mailhog: http://localhost:8025