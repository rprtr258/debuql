# debuql

Библиотека для отладки `sql` запросов.

# :bangbang: Следует использовать только для разработки :bangbang:

## Основные функции

### SDumpQuery
```go
// SDumpQuery - преобразовать запрос с аргументами в одну строку.
func SDumpQuery(query string, args ...any) string
```

Полезно для проверки получившегося запроса и проверки запроса руками.

Пример использования:

```go
import "github.com/rprtr258/debuql"

query, args, err := squirrel.Select(...).From(...).ToSql()
if err != nil { ... }
log.Println(debuql.SDumpQuery(query, args...))
```

### DumpQueryTypes
```go
// DumpQueryTypes - получить типы данных для столбцов указанного запроса.
// Требует доступ с правами на запись, чтобы создать временную таблицу `temp`.
func DumpQueryTypes(db *sqlx.DB, query string, args ...any) Columns
```

Полезно для проверки типов полей (название полей - `field`, тип - `type`) и могут ли они являться nullable (`nullable`).

Пример использования:

```go
import "github.com/rprtr258/debuql"

query, args, err := squirrel.Select(...).From(...).ToSql()
if err != nil { ... }
spew.Dump(debuql.DumpQueryTypes(repo.db, query, args...))
```

Вывод выглядит так:
```php
(debuql.Columns) (len=33 cap=33) field="houses_cards_id" type="bigint(10) unsigned" nullable=false key= default=0 extra=
field="address" type="bigint(20) unsigned" nullable=false key= default=0 extra=
field="full_flats_number" type="varchar(70)" nullable=true key= default=<nil> extra=
field="accounts_id" type="bigint(20) unsigned" nullable=false key= default=0 extra=
field="counters_id" type="bigint(20) unsigned" nullable=false key= default=0 extra=
field="c_serial" type="varchar(32)" nullable=true key= default=<nil> extra=
field="c_use_coefficient" type="tinyint(1)" nullable=false key= default=0 extra=
field="c_max_value" type="int(10) unsigned" nullable=false key= default=10000 extra=
field="c_stamp_date" type="date" nullable=false key= default=<nil> extra=
field="c_check_date" type="date" nullable=true key= default=<nil> extra=
```
