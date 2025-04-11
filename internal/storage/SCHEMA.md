# Production-ready архитектура базы данных (нормализованная) для хранения сообщений Telegram

## Основные сущности

- **chats** — информация о чатах/группах/супергруппах
- **users** — информация об авторах сообщений
- **messages** — сообщения, ссылающиеся на чаты и пользователей

---

## DDL (SQLite, нормализованная)

```sql
CREATE TABLE chats (
    id INTEGER PRIMARY KEY,         -- Telegram chat/group/supergroup ID
    title TEXT NOT NULL,
    type TEXT NOT NULL              -- group, supergroup, private, etc.
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY,         -- Telegram user ID
    username TEXT,
    display_name TEXT
);

CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id INTEGER NOT NULL,           -- FK -> chats.id
    message_id INTEGER NOT NULL,        -- Telegram message ID (уникален в рамках чата)
    author_id INTEGER,                  -- FK -> users.id
    text TEXT,                          -- Текст сообщения
    timestamp INTEGER NOT NULL,         -- unixtime (UTC)
    reply_to_message_id INTEGER,        -- FK -> messages.message_id (в этом же чате)
    FOREIGN KEY(chat_id) REFERENCES chats(id),
    FOREIGN KEY(author_id) REFERENCES users(id),
    UNIQUE(chat_id, message_id)
);

CREATE INDEX idx_messages_chat_time ON messages(chat_id, timestamp);
CREATE INDEX idx_messages_author ON messages(author_id);
```

---

## Описание связей

- `messages.chat_id` → `chats.id`
- `messages.author_id` → `users.id`
- `messages.reply_to_message_id` → `messages.message_id` (в рамках одного чата)

---

## Пример сценария записи

1. При получении нового сообщения:
    - Если chat_id не найден в chats — добавить запись.
    - Если author_id не найден в users — добавить запись.
    - Вставить сообщение в messages (с upsert по chat_id+message_id).

---

## Пример запроса для получения последнего сообщения по чату

```sql
SELECT MAX(timestamp) FROM messages WHERE chat_id = ?;
```

---

## Пример запроса для вставки нового сообщения

```sql
INSERT OR IGNORE INTO messages (chat_id, message_id, author_id, text, timestamp, reply_to_message_id)
VALUES (?, ?, ?, ?, ?, ?);
```

---

## Примечания

- Диапазон выгрузки (M) — настраиваемый, по умолчанию сутки.
- Повторная загрузка за прошлое (backfill) не требуется.
- Структура легко расширяется для хранения media, forwarded, reactions и других метаданных.