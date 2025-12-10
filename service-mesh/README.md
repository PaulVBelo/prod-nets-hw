# Service Mesh: Notes & Mailer

Проект реализует service mesh из двух сервисов **note-service** и **mailer-service**.
Взаимодействие между ними происходит по HTTPS через sidecar-proxy (nginx).

Для эмуляции реальной почты используется MailHog: поднимаем его как локальный SMTP-сервер, принимающий почту на порту 1025 и отображающий письма в своём веб-интерфейсе.

Архитектура включает:
- центральный входной балансировщик
- sidecar-proxy для каждого сервиса

- (всё ещё) шардированное хранилище PostgreSQL
- (всё ещё) поддерживает REST/SOAP/gRPC
- (всё ещё) ... да всё тут есть, я просто скопировал файлы из 2-ого дз.

## Запуск через docker compose:

Желательно иметь нужные образы заранее:
```
docker pull alpine:3.20
docker pull golang:1.24
```

Запуск: 
```
docker compose up --build
```
---

Доступные порты:

- Edge Gateway HTTPS: **https://localhost:8443**
- MailHog UI: **http://localhost:8025**

Полный список ендпоинтов в README самих сервисов.

---

## Пример пайплайна

### Создание note
```bash
curl -k -X POST https://localhost:8443/notes   -H "Content-Type: application/json"   -d '{"title":"Example","description":"Mesh test"}'
```

### GET note by id
```bash
curl -k https://localhost:8443/notes/<NOTE_ID>
```

### Отправка note в письме (Mailhog)
```bash
curl -k -X POST https://localhost:8443/mailer/send   -H "Content-Type: application/json"   -d '{"note_id":"<NOTE_ID>", "email":"user@example.com"}'
```

### 3.4 Посмотреть отправленное письмо

Открыть Mailhog в браузере
```
http://localhost:8025
```