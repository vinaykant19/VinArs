package database

import (
	cnf "../../../configuration"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Database struct {
	DbConn *sql.DB
	Txn []*sql.Tx
}

func (db *Database) Connect (conf *cnf.Configuration) (error){
	connectionString := conf.DB_User + ":" + conf.DB_Pwd + "@tcp(" + conf.DB_Host
	if len(conf.DB_Port) > 0 {
		connectionString += ":" + conf.DB_Port
	}
	connectionString += ")/" + conf.DB_Name

	conn, err := sql.Open(conf.DB_Driver, connectionString)
	if err != nil {
		db.Disconnect()
		//log.Fatal(err)
	}
	db.DbConn = conn

	return err
}

func (db *Database) Disconnect () {
	defer db.DbConn.Close()
}

func (db *Database) IsLive() bool {
	if db.DbConn != nil {
		return true
	}

	return false
}

func (db *Database) Transaction(action string) error {
	if action == "start" {
		txn, err := db.DbConn.Begin()
		db.Txn = append(db.Txn, txn)
		return err
	}
	if action == "commit" {
		lastTxnId := len(db.Txn) - 1
		if lastTxnId >= 0 {
			txn := db.Txn[lastTxnId]
			err := txn.Commit()
			newTxnId := len(db.Txn) - 1
			if newTxnId > 0 {
				db.Txn = db.Txn[0 : newTxnId-1]
			} else if newTxnId == 0 {
				db.Txn = db.Txn[0:0]
			} else {
				db.Txn = nil
			}
			return err
		}
	}
	if action == "rollback" {
		lastTxnId := len(db.Txn) - 1
		if lastTxnId >= 0 {
			txn := db.Txn[lastTxnId]
			err :=txn.Rollback()
			newTxnId := len(db.Txn) - 1
			if newTxnId > 0 {
				db.Txn = db.Txn[0 : newTxnId-1]
			} else if newTxnId == 0 {
				db.Txn = db.Txn[0:0]
			} else {
				db.Txn = nil
			}
			return err
		}
	}

	return nil
}

func (db *Database) Insert(table string, data map[string]string) (int64, error) {
	var args []interface{}
	sql := "INSERT INTO " + table + " ("
	fields := ""
	values := " VALUES ("

	if len(data) > 0 {
		i := 1
		for key, value := range data {
			args = append(args, value)
			if i < len(data) {
				fields += key + ", "
				values += "?, "
			} else {
				fields += key + ")"
				values += "? )"
			}
			i++
		}
		sql += fields + values
		rs, err := db.Exec(sql, args...)
		if err == nil {
			return rs.LastInsertId()
		}

		return 0, err
	}

	return 0, errors.New("empty dataset provided")
}

func (db *Database) Update(table string, toUpdate map[string]string, condition map[string]string) (int64, error) {
	var args []interface{}
	sql := "UPDATE " + table + " SET "

	if len(toUpdate) > 0 {
		i := 1
		for key, value := range toUpdate {
			args = append(args, value)
			if i < len(toUpdate) {
				sql += key +  "=?, "
			} else {
				sql += key +  "=?"
			}
			i++
		}
		if len(condition) > 0 {
			sql += " WHERE "
			i := 1
			for key, value := range condition {
				args = append(args, value)
				if i < len(condition) {
					sql += key +  "=? AND "
				} else {
					sql += key +  "=?"
				}
				i++
			}
		}

		rs, err := db.Exec(sql, args...)
		if err == nil {
			return rs.RowsAffected()
		}

		return 0, err
	}

	return 0, errors.New("empty dataset provided")
}

func (db *Database) Delete(table string, condition map[string]string)(int64, error) {
	var args []interface{}
	sql := "DELETE FROM " + table
	if len(condition) > 0 {
		sql += " WHERE "
		i := 1
		for key, value := range condition {
			args = append(args, value)
			if i < len(condition) {
				sql += key +  "=?, "
			} else {
				sql += key +  "=?"
			}
			i++
		}

		rs, err := db.Exec(sql, args...)
		if err == nil {
			return rs.RowsAffected()
		}
		return 0, err
	}

	return 0, errors.New("empty condition provided")
}

func (db *Database) Select(table string, fields []string, condition map[string]string, orderBy []string, totalRows int) (*sql.Rows, error) {
	var args []interface{}
	sql := "SELECT "

	if len(fields) > 0 {
		for key, value := range fields {
			if key < len(fields) -1 {
				sql += value +  ", "
			} else {
				sql += value
			}
		}
	} else {
		sql += " * "
	}
	sql += " FROM " + table
	if len(condition) > 0 {
		sql += " WHERE "
		i := 1
		for key, value := range condition {
			args = append(args, value)
			if i < len(condition) {
				sql += key +  "=? AND "
			} else {
				sql += key +  "=?"
			}
			i++
		}
	}
	if len(orderBy) > 0 {
		sql += " ORDER BY "
		for key, value := range orderBy {
			if key < len(orderBy) -1 {
				sql += value +  ", "
			} else {
				sql += value
			}
		}
	}

	if (totalRows > 0) {
		sql += " limit 0, " + strconv.Itoa(totalRows)
	}

	rows, err := db.Query(sql, args...)

	if err == nil {
		return rows, nil
	}
	db.Disconnect()
	return nil, err
}

func (db *Database) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	txn := db.LatestTxn()
	if txn != nil {
		rows, err = txn.Query(query, args...)
	} else {
		rows, err = db.DbConn.Query(query, args...)
	}

	if err != nil {
		fmt.Println(err)
	}
	return rows, err
}

func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(query)
	if err == nil {
		rs, err := stmt.Exec(args...)
		return rs, err
	}

	return nil, err
}

func (db *Database) Prepare(query string) (stmt *sql.Stmt, err error) {
	txn := db.LatestTxn()
	if txn != nil {
		stmt, err = txn.Prepare(query)
	} else{
		stmt, err = db.DbConn.Prepare(query)
	}

	return stmt, err
}

func (db *Database) LatestTxn() *sql.Tx {
	lastTxnId := len(db.Txn) - 1
	if lastTxnId >= 0 {
		return db.Txn[lastTxnId]
	}
	return nil
}
