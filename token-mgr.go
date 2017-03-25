package tokenmanager



import (
	"database/sql"
	"github.com/randomtask1155/dbhandler"
	"github.com/gorilla/sessions"
	"fmt"
	"errors"
	"golang.org/x/oauth2"
	"net/http"
)
var(
	// TokenTable table used to store userinfo
	TokenTable = "users"
	// DBType used in dbhandler to find the correct driver
	DBType = "mysql"
)

// TokenTuple defines a single row from the users table
type TokenTuple struct {
	UserID int `json:"userid"`
	UserName string `json:"username"`
	RefreshToken string `json:"refreshtoken"`
}

// ErrNoTokenFound returned when refresh token can not be found
var ErrNoTokenFound = errors.New("Token Manager: No Refresh Token Found")

// CreateSchema will create the users table if it does not exits
func CreateSchema() error {
	dbi, err := dbhandler.NewDBI(DBType) 
	if err != nil {
		return err
	}
	
	// check tables exists 
	_, err = dbi.GetIntValue(fmt.Sprintf("SELECT count(*) from %s", TokenTable))
	if err !=nil {
		// try to create table 
		_, createErr := dbi.SQLSession.Exec("CREATE TABLE users (id serial, username varchar(255), refreshtoken varchar(4096))")
		return createErr
	}
	return nil
}

// GetToken fetches the user refresh token from database
func GetToken(u string) (*TokenTuple, error) {

	t := new(TokenTuple)
	dbi, err := dbhandler.NewDBI(DBType)
	if err != nil {
		return t, err
	}
	defer dbi.Close()


	q := fmt.Sprintf("SELECT * FROM %s WHERE username='%s' limit 1", TokenTable, u)
	err = dbi.SQLSession.QueryRow(q).Scan( &t.UserID, &t.UserName, &t.RefreshToken )
	if err == sql.ErrNoRows {
		return t, ErrNoTokenFound
	} else if err != nil {
		fmt.Printf("GetToken() Error when running query %s: %s\n", q, err)
		return t, err
	}
	return t, nil
}

// UpdateToken inserts new or updates existing users refresh token
func (t *TokenTuple) UpdateToken() error {
	var q string
	var dberr error

	dbi, err := dbhandler.NewDBI(DBType)
	if err != nil {
		return err
	}
	defer dbi.Close()

	tNew, err := GetToken(t.UserName)
	if err == ErrNoTokenFound {
		q = fmt.Sprintf("INSERT INTO %s values(DEFAULT,?,?)", TokenTable)
		_, dberr = dbi.SQLSession.Exec(q, t.UserName, t.RefreshToken)
	}else {
		q = fmt.Sprintf("UPDATE %s SET username='%s', refreshtoken='%s' WHERE id = %d", TokenTable, t.UserName, t.RefreshToken, tNew.UserID)
		_, dberr = dbi.SQLSession.Exec(q)
	}

	if dberr != nil {
		fmt.Printf("UpdateToken() Error when running query %s: %s\n", q, dberr)
		return dberr
	}
	return nil
}

// GetFullToken Given a cookie store we retrieve the access token and fetch the refresh
// token form the database
func GetFullToken(r *http.Request, store *sessions.CookieStore, sessionName string) (*oauth2.Token, error) {
  session, err := store.Get(r, sessionName)
  if err != nil {
    return &oauth2.Token{}, err 
  }
  tok, tokOK := session.Values["AuthToken"].(*oauth2.Token)
  if ! tokOK {
    return tok, fmt.Errorf("Could not find authtoken in cookie store")
  }
	email, emailOK := session.Values["Email"].(string)
	if ! emailOK {
		return tok, fmt.Errorf("Could not find email address in cookie store")
	}
  dbTok, err := GetToken(email)
  if err != nil {
    return tok, err
  }
  tok.RefreshToken = dbTok.RefreshToken
  return tok, nil
}
