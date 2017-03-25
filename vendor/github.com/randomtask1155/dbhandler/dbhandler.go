
package dbhandler

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq" // postgres
	_ "github.com/go-sql-driver/mysql" // mysql
	"encoding/json"
	"os"
	"fmt"
)

var (
	
	// DBTYPE can be mysql or postgres but only mysql implemented for nwo
	//DBTYPE = "mysql" // [postgres || mysql]
	
	// DBURL is the connection string used to open sql session
	//DBURL string
)

// DBInstance manages the sql session
type DBInstance struct {
	SQLSession *sql.DB
	MysqlEnv VCAPServicesMySQL
	DBTYPE string
	DBURL string
}

// NewDBI creates a new DBInstance struct and return it to caller
//	caller is expected to close the database instance dbi.Close()
func NewDBI(dbtype string) (*DBInstance, error) {
	dbi := new(DBInstance)
	dbi.DBTYPE = dbtype 
	dbi.parseEnv()
	err := dbi.ConnectDB()
	if err != nil {
		return nil, err
	}
	return dbi, nil
}

// newMockDBI used for testing purposes 
func newMockDBI(dbtype, dburl string)  (*DBInstance, error) {
	dbi := new(DBInstance)
	dbi.DBTYPE = dbtype 
	dbi.DBURL = dburl
	err := dbi.ConnectDB()
	if err != nil {
		return nil, err
	}
	return dbi, nil
}

// ConnectDB creates a new database session.  Caller needs to call Close() when done
func (dbi *DBInstance) ConnectDB() error {
	sess, err := sql.Open(dbi.DBTYPE, dbi.DBURL)
	if err != nil {
		return errors.New("can not connect to database: " + err.Error())
	}
	dbi.SQLSession = sess
	dbi.SQLSession.SetMaxOpenConns(1) // make sure there is only one session open with database at a time
	return nil
}

// parseEnv for vcap services
// TODO add postgres support
func (dbi *DBInstance)parseEnv() {
	switch {
	case dbi.DBTYPE == "mysql":
		VCAP := VCAPServicesMySQL{}
		setEnv(&VCAP)
		if len(VCAP.MySQL) > 0{
			dbi.DBURL = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			 													VCAP.MySQL[0].Credentials.Username,
																VCAP.MySQL[0].Credentials.Password,
																VCAP.MySQL[0].Credentials.Hostname,
																VCAP.MySQL[0].Credentials.Port,
																VCAP.MySQL[0].Credentials.Name)
		}
	}
}

// parse VCAP_SERVICES environment into struct pointer 
func setEnv(d interface{}) {
	VCAP := os.Getenv("VCAP_SERVICES")
	if VCAP == "" {
		return // no environment found so use whatever DBURL is set to
	}
	b := []byte(VCAP)
	err := json.Unmarshal(b, d) 
	if err != nil {
		fmt.Printf("dbhandler:setEnv:ERROR:%s", err)
	}
}

// Close the database session
func (dbi *DBInstance) Close() error {
	if dbi.SQLSession != nil {
		err := dbi.SQLSession.Close()
		if err != nil {
			return errors.New("can not close libpq session: " + err.Error())
		}
	}
	return nil
}

/*####################################################################*/
/*
	QUERY FUNCTIONS
*/

// GetRowSet returns row set from query and expects caller to handle error like sql.ErrNoRows
func (dbi *DBInstance) GetRowSet(qstring string) (*sql.Rows, error) {
	return dbi.SQLSession.Query(qstring)

}

// GetStringList returns string slice from rowset
func (dbi *DBInstance) GetStringList(qstring string) ([]string, error) {
	var result []string
	result = make([]string, 0)
	rows, err := dbi.GetRowSet(qstring)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return result, err
		}
		result = append(result, s)
	}
	return result, nil
}

// GetIntList returns int slice from rowset
func (dbi *DBInstance) GetIntList(qstring string) ([]int, error) {
	var result []int 
	result = make([]int, 0)
	rows, err := dbi.GetRowSet(qstring)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var s int
		err := rows.Scan(&s)
		if err != nil {
			return result, err
		}
		result = append(result, s)
	}
	return result, nil
}

// GetIntValue Assumed single row/column query result of integer type and returns that value
func (dbi *DBInstance) GetIntValue(qstring string) (int, error) {
	var v int
	err := dbi.SQLSession.QueryRow(qstring).Scan(&v)
	if err != nil {
		return v, err
	}
	return v, nil
}

// GetStringValue Assumed single row/column query result of string type and returns that value
func (dbi *DBInstance) GetStringValue(qstring string) (string, error) {
	var v string
	err := dbi.SQLSession.QueryRow(qstring).Scan(&v)
	if err != nil {
		return v, err
	}
	return v, nil
}
