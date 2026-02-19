# pageable

A Go library for REST API pagination. Supports offset-based and cursor-based pagination with generic responses. Zero dependencies.

## Install

```bash
go get github.com/ishinvin/pageable
```

Requires Go 1.21+

## Offset-Based Pagination

```go
func listUsers(w http.ResponseWriter, r *http.Request) {
    // Parse from query string: ?page=2&size=20&sort=name,desc&sort=id,asc
    req := pageable.PageRequestFromQuery(r.URL.Query()).
        SortableFields("id", "name", "created_at").
        WithDefaultSort(pageable.Sort{Field: "id", Direction: pageable.ASC})

    users, total := queryUsers(req.Offset(), req.Limit(), req.OrderBy())

    page := pageable.NewPage(users, req, total)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(page)
}
```

```json
{
  "items": [
    { "id": 21, "name": "Alice" },
    { "id": 22, "name": "Bob" }
  ],
  "metadata": {
    "page": 2,
    "size": 20,
    "totalItems": 95,
    "totalPages": 5
  }
}
```

## Cursor-Based Pagination

```go
func listPosts(w http.ResponseWriter, r *http.Request) {
    req := pageable.CursorRequestFromQuery(r.URL.Query()).
        SortableFields("id", "created_at").
        WithDefaultSort(pageable.Sort{Field: "id", Direction: pageable.ASC})

    var posts []Post
    if req.HasCursor() {
        cursorData, err := req.DecodedCursor()
        if err != nil {
            http.Error(w, "invalid cursor", http.StatusBadRequest)
            return
        }
        switch cursorData.Direction {
        case pageable.Prev:
            posts = queryPostsBefore(cursorData.Value, req.Limit(), req.OrderBy())
        default:
            posts = queryPostsAfter(cursorData.Value, req.Limit(), req.OrderBy())
        }
    } else {
        posts = queryPosts(req.Limit(), req.OrderBy())
    }

    // Detect hasNext by fetching Size+1 rows via req.Limit()
    hasNext := len(posts) > req.Size
    if hasNext {
        posts = posts[:req.Size]
    }

    var nextCursor, prevCursor string
    if len(posts) > 0 {
        if hasNext {
            last := posts[len(posts)-1]
            nextCursor, _ = pageable.EncodeCursor(pageable.CursorData{
                Value:     fmt.Sprintf("%d", last.ID),
                Direction: pageable.Next,
            })
        }
        if req.HasCursor() {
            first := posts[0]
            prevCursor, _ = pageable.EncodeCursor(pageable.CursorData{
                Value:     fmt.Sprintf("%d", first.ID),
                Direction: pageable.Prev,
            })
        }
    }

    page := pageable.NewCursorPage(posts, nextCursor, prevCursor, hasNext, req.HasCursor(), req.Size)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(page)
}
```

```json
{
  "items": [
    { "id": 42, "title": "Hello" },
    { "id": 43, "title": "World" }
  ],
  "metadata": {
    "nextCursor": "eyJ2IjoiNDMiLCJkIjoibmV4dCJ9",
    "prevCursor": "eyJ2IjoiNDIiLCJkIjoicHJldiJ9",
    "hasNext": true,
    "hasPrev": true,
    "size": 25
  }
}
```

## Sorting

Sort parameters use `field,direction` format (repeatable):

```
?sort=name,desc&sort=id,asc
```

`SortableFields` whitelists allowed fields to prevent SQL injection. `WithDefaultSort` provides a fallback when no sort is given.

```go
req := pageable.PageRequestFromQuery(r.URL.Query()).
    SortableFields("id", "name", "created_at").
    WithDefaultSort(pageable.Sort{Field: "id", Direction: pageable.ASC})

req.OrderBy() // "name desc, id asc"
```

## Compound Cursors

For cursors that need multiple values (e.g., `created_at` + `id` for stable ordering):

```go
cursor, _ := pageable.EncodeCursor(pageable.CursorData{
    Value:     "42",
    Direction: pageable.Next,
    Extra:     map[string]string{"created_at": "2024-01-15T10:30:00Z"},
})

data, _ := pageable.DecodeCursor(cursor)
// data.Value == "42"
// data.Direction == "next"
// data.Extra["created_at"] == "2024-01-15T10:30:00Z"
```

## Empty Pages

```go
page := pageable.EmptyPage[User](req)            // offset-based
page := pageable.EmptyCursorPage[Post](req.Size)  // cursor-based
```

## Query Parameters

| Parameter | Default | Description |
|:----------|:--------|:------------|
| `page` | 1 | Page number (1-indexed, offset only) |
| `cursor` | — | Encoded cursor token (cursor only) |
| `size` | 10 | Items per page (max 1000) |
| `sort` | — | Sort field: `field,direction` (repeatable) |

## Documentation

Full documentation is available at [ishinvin.github.io/pageable](https://ishinvin.github.io/pageable).

## License

MIT
