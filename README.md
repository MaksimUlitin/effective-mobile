
# Решение тестового задания на позицию Junior Golang Developer для Effective-Mobile

# Запуск:
``` 
git clone https://github.com/MaksimUlitin/effective-mobile.git
```
```bash
make all
```

# Документация по методам API

- [Вызов методов](#вызов-методов)

Методы:

- [Songs](#songs)
    - [Add song information](#add-song-information)
    - [List songs](#list-songs)
    - [List songs with filter](#list-songs-with-filter)
    - [Update an existing song](#update-an-existing-song)
    - [Delete a song](#delete-a-song)
    - [Get song text by ID with pagination](#get-song-text-by-id-with-pagination)

---

## Вызов методов

Методы вызываются через HTTP с использованием методов GET, POST, PUT или DELETE.

Пример вызова метода:

```
http://localhost:8080/<resource>?<params>
```


---

## Songs

### Add song information

Добавляет информацию о песне.

#### URL

```
POST /info
```

#### Параметры

| Параметр | Тип   | Описание             | Обязательный |
|----------|-------|----------------------|--------------|
| group    | string | Название группы      | Да           |
| song     | string | Название песни       | Да           |

#### Пример запроса

```json
{
  "group": "Muse",
  "song": "Supermassive Black Hole"
}
```

#### Пример ответа

```json
{
  "release_date": "16.07.2006",
  "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
  "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
}
```

---

### List songs with filter

Возвращает список песен с фильтрацией.

#### URL

```
GET /songs
```

#### Параметры фильтрации

| Параметр | Тип   | Описание                  | Обязательный |
|----------|-------|---------------------------|--------------|
| page     | int   | Номер страницы            | Нет          |
| limit    | int   | Количество элементов      | Нет          |
| group    | string | Фильтр по названию группы | Нет          |
| song     | string | Фильтр по названию песни  | Нет          |
| link     | string | Фильтр по ссылке песни    | Нет          |

#### Пример запроса

```
GET http://localhost:8080/songs?page=1&limit=10&group=Muse
```

#### Пример ответа

```json
[
  {
    "id": 11,
    "group": "Muse",
    "song": "Supermassive Black Hole",
    "release_date": "16.07.2006",
    "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
    "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
    "created_at": "2024-12-04T02:22:46.724854+03:00",
    "updated_at": "2024-12-04T02:22:46.724854+03:00"
  }
]
```

---

### Update an existing song

Обновляет информацию о существующей песне.

#### URL

```
PUT /songs/{id}
```

#### Параметры

| Параметр     | Тип   | Описание                    | Обязательный |
|--------------|-------|-----------------------------|--------------|
| id (в URL)   | int   | Идентификатор песни         | Да           |
| group        | string | Название группы            | Да           |
| song         | string | Название песни             | Да           |
| release_date | string | Дата выпуска               | Нет          |
| text         | string | Текст песни                | Нет          |
| link         | string | Ссылка на песню            | Нет          |

#### Пример запроса

```json
{
  "group": "Muse",
  "song": "Supermassive Black Hole",
  "release_date": "16.07.2006",
  "text": "sample text",
  "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
}
```

#### Пример ответа

```json
{
  "message": "song updated successfully"
}
```

---

### Delete a song

Удаляет песню по ID.

#### URL

```
DELETE /songs/{id}
```

#### Параметры

| Параметр     | Тип   | Описание            | Обязательный |
|--------------|-------|---------------------|--------------|
| id (в URL)   | int   | Идентификатор песни | Да           |

#### Пример запроса

```
DELETE http://localhost:8080/songs/1
```

#### Пример ответа

```json
{
  "id #1": "deleted"
}
```

---

### Get song text by ID with pagination

Возвращает текст песни с поддержкой пагинации.

#### URL

```
GET /songs/{id}/text
```

#### Параметры

| Параметр     | Тип   | Описание                    | Обязательный |
|--------------|-------|-----------------------------|--------------|
| id (в URL)   | int   | Идентификатор песни         | Да           |
| page         | int   | Номер страницы текста       | Нет          |
| limit        | int   | Количество строк на странице | Нет         |

#### Пример запроса

```
GET http://localhost:8080/songs/1/text?page=1&limit=10
```

#### Пример ответа

```json
{
  "limit": 10,
  "page": 1,
  "songId": 11,
  "text": [
    "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?",
    "Ooh\nYou set my soul alight\nOoh\nYou set my soul alight"
  ],
  "total": 2,
  "totalPage": 1
}
```