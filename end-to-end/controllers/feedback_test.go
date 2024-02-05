package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"homethings.ytsruh.com/lib"
)

func init() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func Test_CreateFeedback(t *testing.T) {
	// Setup mock database & add mock rows
	db, mock, err := lib.SetupMockDB()
	if err != nil {
		t.Errorf(err.Error())
	}
	// Setup test cases
	type args struct {
		Title string
		Body  interface{}
	}
	tests := []struct {
		name   string
		claims *CustomClaims
		args   args
		want   int
	}{
		{
			name: "Feedback created",
			claims: &CustomClaims{
				User:      "testing@gmail.com",
				Id:        "1",
				AccountId: "1234567890",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
					Issuer:    "homethings",
				},
			},
			args: args{
				Title: "Test Title",
				Body:  "Test Body",
			},
			want: 200,
		},
		{
			name: "Bad request",
			claims: &CustomClaims{
				User:      "notfound@gmail.com",
				Id:        "123456",
				AccountId: "1234567890",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
					Issuer:    "homethings",
				},
			},
			args: args{
				Title: "Test Title",
				Body:  123,
			},
			want: 400,
		},
		{
			name: "Failed validation",
			claims: &CustomClaims{
				User:      "notfound@gmail.com",
				Id:        "123456",
				AccountId: "1234567890",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
					Issuer:    "homethings",
				},
			},
			args: args{
				Body: "Example body",
			},
			want: 400,
		},
	}
	// Run test cases
	for _, test := range tests {
		// Setup expected SQL response
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO \"feedback\" \\(\"title\",\"body\",\"user_id\",\"created_at\",\"updated_at\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING \"id\"").
			WithArgs("Test Title", "Test Body", "1", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		// Run test
		t.Run(test.name, func(t *testing.T) {
			// Setup echo
			e := echo.New()
			api := &API{DB: db}
			e.Use(api.SetJWTAuth())
			// Marshal test.args into JSON & convert to io.Reader/buffer
			jsonData, err := json.Marshal(test.args)
			if err != nil {
				t.Fatal(err)
			}
			body := bytes.NewBuffer(jsonData)
			req := httptest.NewRequest("POST", "/v1/feedback", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/v1/feedback")
			// Create user token & set user context
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, test.claims)
			c.Set("user", token)
			// Run
			err = api.CreateFeedback(c)
			if err != nil {
				t.Errorf("gave error : %v", err)
				return
			}
			// Check expected response code
			if rec.Code != test.want {
				t.Errorf("gave code: %v, wanted code: %v", rec.Code, test.want)
			}
			// if err := mock.ExpectationsWereMet(); err != nil {
			// 	t.Errorf("there were unfulfilled expectations: %s", err)
			// }
		})
	}

}
