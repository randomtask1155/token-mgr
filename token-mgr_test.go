package tokenmanager 

import (
  "testing"
  "github.com/gorilla/sessions"
  "golang.org/x/oauth2"
  "net/http"
  "time"
  "fmt"
  "os"
)

var (
  MockTokenTuple TokenTuple
  MockUserEmail = "test@testuser.domain"
)

/*
  Givne database is set in VCAP_SERVICES envrionment we make sure schema
  exists and we add or update a test user to database
*/
func init(){
  
  err := CreateSchema()
  if err != nil {
    fmt.Printf("token-mgr_test:init:%s\n", err)
    os.Exit(1)
  }
  MockTokenTuple = TokenTuple{1234, MockUserEmail, "12345678901234567890"}
  err = MockTokenTuple.UpdateToken()
  if err != nil {
    fmt.Printf("token-mgr_test:init:%s\n", err)
    os.Exit(2)
  }
}

func TestGetFullToken(t *testing.T) {
  mockR := &http.Request{}
  mockStore := sessions.NewCookieStore([]byte("TestGetFullToken"))
  mockSession, err := mockStore.Get(mockR, "mock-session")
  if err != nil {
    t.Fatal(err)
  }
  mockToken := oauth2.Token{}
  mockToken.AccessToken = "somerandomtokenvalue12345" 
  mockToken.TokenType = "Bearer"
  mockToken.RefreshToken = ""
  mockToken.Expiry = time.Now()
  mockSession.Values["Email"] = MockUserEmail
  mockSession.Values["AuthToken"] = &mockToken
  tok, err := GetFullToken(mockR, mockStore, "mock-session")
  if err != nil {
    t.Fatal(err)
  }
  if tok.RefreshToken == "" {
    t.Fatal("No refresh token retrieved")
  }
}