package debuql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// as - преобразовать x в T, если возможно. Возвращает результат конвертации
// и флаг, успешно ли сконвертировалось значение.
func as[T any](x any) (T, bool) {
	var t T
	typeOfT := reflect.TypeOf(t)

	if !reflect.TypeOf(x).ConvertibleTo(typeOfT) {
		return t, false
	}

	return reflect.
		ValueOf(x).
		Convert(typeOfT).
		Interface().(T), true
}

// SDumpQuery - преобразовать запрос с аргументами в одну строку.
// ИСПОЛЬЗОВАТЬ ИСКЛЮЧИТЕЛЬНО ДЛЯ ТЕСТИРОВАНИЯ И ОТЛАДКИ.
func SDumpQuery(query string, args ...any) string {
	for _, arg := range args {
		var argStr string

		if u, ok := as[uint64](arg); ok {
			argStr = strconv.FormatUint(u, 10)
		} else if i, ok := as[int](arg); ok {
			argStr = strconv.Itoa(i)
		} else if t, ok := as[time.Time](arg); ok {
			argStr = t.Format("TIMESTAMP'2006.01.02 15:04:05'")
		} else if s, ok := as[string](arg); ok {
			argStr = fmt.Sprintf(`"%s"`, s)
		} else if s, ok := as[bool](arg); ok {
			if s {
				argStr = "1"
			} else {
				argStr = "0"
			}
		} else {
			err := fmt.Errorf("unknown type to sql-dump: %[1]T=%#[1]v", arg)
			panic(err)
		}

		query = strings.Replace(
			query,
			"?",
			argStr,
			1,
		)
	}
	return query
}

func marshalSQLNullString(s sql.NullString) string {
	if !s.Valid {
		return "<nil>"
	}

	return s.String
}

// ColumnType - описание БД типа столбца.
// См. подробнее здесь https://dev.mysql.com/doc/refman/8.0/en/show-columns.html
type ColumnType struct {
	// Field - название столбца
	Field string
	// Type - тип БД стобца
	Type string
	// Nullable - является ли тип NULLABLE
	Nullable bool
	// Key -  является ли столбец индексируемым
	Key sql.NullString
	// Default - дефолтное значение для столбца
	Default sql.NullString
	// Extra дополнительная информация о столбце
	Extra sql.NullString
}

type Columns []ColumnType

func (cols Columns) String() string {
	sb := strings.Builder{}
	for _, row := range cols {
		sb.WriteString(
			fmt.Sprintf(
				"field=%#v type=%#v nullable=%#v key=%s default=%s extra=%s\n",
				row.Field,
				row.Type,
				row.Nullable,
				marshalSQLNullString(row.Key),
				marshalSQLNullString(row.Default),
				marshalSQLNullString(row.Extra),
			),
		)
	}
	return sb.String()
}

// DumpQueryTypes - получить типы данных для столбцов указанного запроса.
// ИСПОЛЬЗОВАТЬ ИСКЛЮЧИТЕЛЬНО ДЛЯ ТЕСТИРОВАНИЯ И ОТЛАДКИ.
// Требует доступ с правами на запись, чтобы создать временную таблицу `temp`.
func DumpQueryTypes(db *sqlx.DB, query string, args ...any) Columns {
	if _, err := db.Exec("DROP TABLE IF EXISTS `temp`"); err != nil {
		panic(err.Error())
	}

	if _, err := db.Exec(
		"CREATE TEMPORARY TABLE `temp` "+strings.Replace(query, ";", "", -1)+" LIMIT 0",
		args...,
	); err != nil {
		panic(err.Error())
	}

	type destRow struct {
		Field   string         `db:"Field"`
		Type    string         `db:"Type"`
		Null    string         `db:"Null"`
		Key     sql.NullString `db:"Key"`
		Default sql.NullString `db:"Default"`
		Extra   sql.NullString `db:"Extra"`
	}
	dest := []destRow{}
	if err := db.Select(&dest, "DESCRIBE `temp`;"); err != nil {
		panic(err.Error())
	}

	res := make([]ColumnType, len(dest))
	for i, row := range dest {
		res[i] = ColumnType{
			Field:    row.Field,
			Type:     row.Type,
			Nullable: row.Null == "YES",
			Key:      row.Key,
			Default:  row.Default,
			Extra:    row.Extra,
		}
	}
	return res
}
