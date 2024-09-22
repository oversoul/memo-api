# Memo

Memo is an api for saving notes.

it focuses on having different templates based on the type of note.

- todo
- movie (or shows to watch)
- note

the idea is to have a separate way to display said notes, and also makes it easier to navigate.


# Data-Layout

Memo uses MongoDB to store the data.

- Sample document for a general note

```
{
    "_id": ObjectId("5f8a7b2e1c9d440000a1e345"),
    "type": "text",
    "title": "Project Ideas",
    "text_note": {
        "content": "1. Develop a personal finance app\n2. Create a recipe sharing platform"
    }
    "tags": ["ideas", "projects"],
    "created_at": ISODate("2023-09-15T10:30:00Z"),
    "updated_at": ISODate("2023-09-15T10:30:00Z")
}
```

- Sample document for a todo item

```
{
    "_id": ObjectId("5f8a7b2e1c9d440000a1e346"),
    "type": "todo",
    "title": "Buy groceries",
    "todo_note": {
        "tasks": [
            {
                _id: ObjectId('66ef669a910d2b2dea78ffd8'),
                content: 'Send Create Memo form',
                is_completed: false,
                completed_at: null
            },
        ]
    },
    "tags": ["shopping", "personal"],
    "created_at": ISODate("2023-09-15T11:00:00Z"),
    "updated_at": ISODate("2023-09-15T11:00:00Z")
}
```

- Sample document for a movie to watch

```
{
    "_id": ObjectId("5f8a7b2e1c9d440000a1e347"),
    "type": "movie",
    "title": "Inception",
    "movie_note": {
        "director": "Christopher Nolan",
        "year": 2010,
        "watched": false,
    },
    "tags": ["must-watch", "recommended"],
    "created_at": ISODate("2023-09-15T12:15:00Z"),
    "updated_at": ISODate("2023-09-15T12:15:00Z")
}
```
