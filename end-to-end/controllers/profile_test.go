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

func Test_GetProfile(t *testing.T) {
	// Setup mock database & add mock rows
	db, mock, err := lib.SetupMockDB()
	if err != nil {
		t.Errorf(err.Error())
	}
	rows := mock.NewRows([]string{"id", "name", "email", "profile_image", "show_books", "show_documents"}).
		AddRow("1", "testing", "testing@gmail.com", "image.jpg", true, true)
	// Setup test cases
	tests := []struct {
		name   string
		claims *CustomClaims
		want   int
	}{
		{
			name: "Profile found",
			claims: &CustomClaims{
				User:      "testing@gmail.com",
				Id:        "1",
				AccountId: "1234567890",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
					Issuer:    "homethings",
				},
			},
			want: 200,
		},
		{
			name: "Profile not found",
			claims: &CustomClaims{
				User:      "notfound@gmail.com",
				Id:        "123456",
				AccountId: "1234567890",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
					Issuer:    "homethings",
				},
			},
			want: 500,
		},
	}
	// Run test cases
	for _, test := range tests {
		// Setup expected SQL response
		mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE id = \\$1 ORDER BY \"users\"\\.\"id\" LIMIT 1$").
			WithArgs(test.claims.Id).
			WillReturnRows(rows)
		// Run test
		t.Run(test.name, func(t *testing.T) {
			// Setup echo
			e := echo.New()
			api := &API{DB: db}
			e.Use(api.SetJWTAuth())
			req := httptest.NewRequest("GET", "/v1/profile", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/v1/profile")
			// Create user token & set user context
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, test.claims)
			c.Set("user", token)
			// Run
			err = api.GetProfile(c)
			if err != nil {
				t.Errorf("gave error : %v", err)
				return
			}
			// Check expected response code
			if rec.Code != test.want {
				t.Errorf("gave code: %v, wanted code: %v", rec.Code, test.want)
			}
		})
	}

}

func Test_PatchProfile(t *testing.T) {
	// Setup mock database & add mock rows
	db, mock, err := lib.SetupMockDB()
	if err != nil {
		t.Errorf(err.Error())
	}
	rows := mock.NewRows([]string{"id", "name", "email", "profile_image", "show_books", "show_documents"}).
		AddRow("1", "testing", "testing@gmail.com", "image.jpg", true, true)
		// Setup test cases
	type args struct {
		Name          string
		ProfileImage  string
		ShowBooks     bool
		ShowDocuments bool
	}
	tests := []struct {
		name   string
		claims *CustomClaims
		args   args
		want   int
	}{
		{
			name: "Update Success",
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
				Name:          "newname",
				ProfileImage:  "test.png",
				ShowBooks:     false,
				ShowDocuments: false,
			},
			want: 200,
		},
	}
	// Run test cases
	for _, test := range tests {
		// Setup expected SQL response
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE \"users\" SET \"name\"=\\$1,\"profile_image\"=\\$2,\"updated_at\"=\\$3 WHERE id = \\$4").
			WithArgs("newname", "test.png", sqlmock.AnyArg(), test.claims.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE id = \\$1 ORDER BY \"users\"\\.\"id\" LIMIT 1$").
			WithArgs(test.claims.Id).
			WillReturnRows(rows)
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
			req := httptest.NewRequest("PATCH", "/v1/profile", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/v1/profile")
			// Create user token & set user context
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, test.claims)
			c.Set("user", token)
			// Run
			err = api.PatchProfile(c)
			if err != nil {
				t.Errorf("gave error : %v", err)
				return
			}
			// Check expected response code
			if rec.Code != test.want {
				t.Errorf("gave code: %v, wanted code: %v", rec.Code, test.want)
			}
		})
	}

}
