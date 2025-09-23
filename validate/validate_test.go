package validate

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidateAndParseQuery(t *testing.T) {
	type User struct {
		Name  string `query:"name" default:"Mark"`
		Age   int    `query:"age" validate:"required"`
		Email string `query:"email" validate:"required,email"`
		Role  string `query:"role" validate:"oneof=admin user"`
	}
	t.Run("valid query parameterss", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/?name=John&age=30&email=john@example.com&role=admin", nil)
		c.Request = req
		var user User
		err := ValidateAndParseQuery(c, &user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("missing required query parameter", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/?age=30", nil)
		c.Request = req
		var user User
		err := ValidateAndParseQuery(c, &user)
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})
	t.Run("invalid email format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/?name=John&age=30&email=john&role=admin", nil)
		c.Request = req
		var user User
		err := ValidateAndParseQuery(c, &user)
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})
	t.Run("invalid role value", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/?name=John&age=30&email=john@example.com&role=guest", nil)
		c.Request = req
		var user User
		err := ValidateAndParseQuery(c, &user)
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	type PageInfo struct {
		Page     int `query:"page" default:"1"`
		PageSize int `query:"limit" default:"10"`
	}
	t.Run("default values", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/?page=2", nil)
		c.Request = req
		var page PageInfo
		err := ValidateAndParseQuery(c, &page)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if page.Page != 2 || page.PageSize != 10 {
			t.Errorf("unexpected page info: got %+v", page)
		}
	})
}

func BenchmarkValidateAndParseQuery(b *testing.B) {
	type User struct {
		Name string `query:"name" validate:"required"`
		Age  int    `query:"age" validate:"required"`
	}
	b.Run("valid query parameterss", func(b *testing.B) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/?name=John&age=30", nil)
		c.Request = req
		var user User
		for b.Loop() {
			ValidateAndParseQuery(c, &user)
		}
	})
	b.Run("missing required query parameter", func(b *testing.B) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/?age=30", nil)
		c.Request = req
		var user User
		for b.Loop() {
			ValidateAndParseQuery(c, &user)
		}
	})
}
func TestValidateAndParseJSON(t *testing.T) {
	type User struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"required"`
	}

	t.Run("valid JSON body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBody := `{"name":"John","age":30}`
		c.Request = httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		var user User
		err := ValidateAndParseJSON(c, &user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user.Name != "John" || user.Age != 30 {
			t.Errorf("unexpected user data: got %+v", user)
		}
	})

	t.Run("missing required field in JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBody := `{"name":"John"}`
		c.Request = httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		var user User
		err := ValidateAndParseJSON(c, &user)
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("invalid JSON format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBody := `{"name":"John","age":30` // Missing closing brace
		c.Request = httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		var user User
		err := ValidateAndParseJSON(c, &user)
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("empty JSON body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", strings.NewReader(""))
		c.Request.Header.Set("Content-Type", "application/json")

		var user User
		err := ValidateAndParseJSON(c, &user)
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})
}
func BenchmarkValidateAndParseJSON(b *testing.B) {
	type User struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"required"`
	}

	b.Run("valid JSON body", func(b *testing.B) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBody := `{"name":"John","age":30}`
		c.Request = httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		var user User
		for b.Loop() {
			ValidateAndParseJSON(c, &user)
		}
	})

	b.Run("missing required field in JSON", func(b *testing.B) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBody := `{"name":"John"}`
		c.Request = httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		var user User
		for b.Loop() {
			ValidateAndParseJSON(c, &user)
		}
	})
}
