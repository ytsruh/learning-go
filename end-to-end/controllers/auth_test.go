package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"ytsruh.com/endtoend/lib"
)

func init() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func Test_Login(t *testing.T) {
	// Setup mock database
	db, mock, err := lib.SetupMockDB()
	if err != nil {
		t.Errorf(err.Error())
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte("testing"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "testing@gmail.com", string(bytes))
	// Setup mock expectations
	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE email = \\$1 ORDER BY \"users\".\"id\" LIMIT 1$").WithArgs("testing@gmail.com").WillReturnRows(rows)
	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE email = \\$1 ORDER BY \"users\".\"id\" LIMIT 1$").WithArgs("testing@gmail.com").WillReturnError(sql.ErrNoRows)
	// Setup test cases
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Successful login",
			args: args{
				email:    "testing@gmail.com",
				password: "testing",
			},
			want: 200,
		},
		{
			name: "Unauthorised login",
			args: args{
				email:    "testing@gmail.com",
				password: "wrongpassword",
			},
			want: 401,
		},
		{
			name: "Bad Request",
			args: args{
				email: "testing@gmail.com",
			},
			want: 400,
		},
	}
	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			api := &API{DB: db}
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(fmt.Sprintf(`{"email":"%s","password":"%s"}`, test.args.email, test.args.password)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/v1/login")
			// Assertions
			err := api.Login(c)
			if err != nil {
				t.Errorf("gave error : %v", err)
				return
			}
			if rec.Code != test.want {
				t.Errorf("gave code: %v, wanted code: %v", rec.Code, test.want)
			}
		})
	}

}
