package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/config"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	"github.com/u104rak1/pocgo/internal/domain/mock"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/seed"
	"github.com/u104rak1/pocgo/internal/server"
	"github.com/u104rak1/pocgo/pkg/strutil"
	"github.com/uptrace/bun"
)

type HTTPRequest struct {
	URL    string         `json:"url"`
	Method string         `json:"method"`
	Header http.Header    `json:"header"`
	Body   interface{}    `json:"body"`
	Query  string         `json:"query"`
	Cookie []*http.Cookie `json:"cookie"`
}

type HTTPResponse struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

type TestResult struct {
	BeforeDB map[string]interface{} `json:"beforeDB"`
	AfterDB  map[string]interface{} `json:"afterDB"`
	Request  HTTPRequest            `json:"request"`
	Response HTTPResponse           `json:"response"`
}

func BeforeAll(t *testing.T) (*echo.Echo, *goldie.Goldie, *bun.DB) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	os.Chdir(filepath.Join(basepath, "../.."))

	db, err := config.LoadDB()
	if err != nil {
		t.Fatal(err)
	}
	ClearDB(t, db)

	e := server.SetupEcho(db)

	goldieDir := filepath.Join("test", "integration", "testdata")
	gol := goldie.New(t, goldie.WithFixtureDir(goldieDir), goldie.WithDiffEngine(goldie.ColoredDiff))

	return e, gol, db
}

func AfterAll(t *testing.T, db *bun.DB) {
	config.CloseDB(db)
}

func ClearDB(t *testing.T, db *bun.DB) {
	for i := len(model.Models) - 1; i >= 0; i-- {
		model := model.Models[i]
		modelType := reflect.TypeOf(model).Elem()
		tableName := modelType.Field(0).Tag.Get("bun")[6:] // Remove "table:" from `bun:"table:table_name"`

		if tableName == "" {
			t.Fatalf("could not retrieve table name for model: %v", model)
		}
		_, err := db.NewTruncateTable().Model(model).Cascade().Exec(context.Background())
		if err != nil {
			t.Fatalf("failed to clear table %s: %v", tableName, err)
		}
	}
}

// Insert test data into the database.
// Master data is inserted by default.
// Please specify a structure or a slice of a structure for the argument models.
func InsertTestData(t *testing.T, db *bun.DB, models ...interface{}) {
	seed.InsertMasterData(db)
	for _, model := range models {
		val := reflect.ValueOf(model)
		if val.Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				elem := val.Index(i).Interface()
				_, err := db.NewInsert().Model(elem).Exec(context.Background())
				if err != nil {
					t.Fatalf("failed to insert test data for model %v: %v", elem, err)
				}
			}
		} else {
			_, err := db.NewInsert().Model(model).Exec(context.Background())
			if err != nil {
				t.Fatalf("failed to insert test data for model %v: %v", model, err)
			}
		}
	}
}

func GetDBData(t *testing.T, db *bun.DB, usedTables []string) map[string]interface{} {
	data := make(map[string]interface{})

	for _, table := range usedTables {
		var records []map[string]interface{}

		err := db.NewSelect().Table(table).Scan(context.Background(), &records)
		if err != nil {
			t.Fatalf("failed to retrieve data from table %s: %v", table, err)
		}

		// Convert []uint8 to string for id fields
		idPattern := regexp.MustCompile(`(^id$|.*_id$)`)
		for _, record := range records {
			for field, value := range record {
				if idPattern.MatchString(field) {
					if id, ok := value.([]uint8); ok {
						record[field] = string(id)
					}
				}
			}
		}

		data[table] = records
	}

	return data
}

func SetAccessToken(t *testing.T, userID string, req *http.Request) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock.NewMockIAuthenticationRepository(ctrl)
	mockUserRepo := mock.NewMockIUserRepository(ctrl)
	service := authDomain.NewService(mockAuthRepo, mockUserRepo)

	ctx := context.Background()
	env := config.NewEnv()
	jwtSecretKey := []byte(env.JWT_SECRET_KEY)
	token, err := service.GenerateAccessToken(ctx, userID, jwtSecretKey)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
}

func NewJSONRequest(t *testing.T, method, url string, requestBody interface{}) (*http.Request, *httptest.ResponseRecorder) {
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	return req, rec
}

func GenerateResultJSON(
	t *testing.T,
	beforeDBData,
	afterDBData map[string]interface{},
	req *http.Request,
	rec *httptest.ResponseRecorder,
	requestBody interface{},
) []byte {
	var responseBody interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	result := TestResult{
		BeforeDB: beforeDBData,
		AfterDB:  afterDBData,
		Request: HTTPRequest{
			URL:    req.URL.String(),
			Method: req.Method,
			Header: req.Header,
			Body:   requestBody,
			Query:  req.URL.RawQuery,
			Cookie: req.Cookies(),
		},
		Response: HTTPResponse{
			StatusCode: rec.Code,
			Body:       responseBody,
		},
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal result to JSON: %v", err)
	}
	return resultJSON
}

func ReplaceDynamicValue(jsonData []byte, replaceKeys []string) []byte {
	for _, key := range replaceKeys {
		camelCasePattern := regexp.MustCompile(`"` + key + `":\s*".*?"`)
		snakeCaseKey := strutil.ToSnakeFromCamel(key)
		snakeCasePattern := regexp.MustCompile(`"` + snakeCaseKey + `":\s*".*?"`)

		jsonData = camelCasePattern.ReplaceAll(jsonData, []byte(`"`+key+`": "ANY"`))
		jsonData = snakeCasePattern.ReplaceAll(jsonData, []byte(`"`+snakeCaseKey+`": "ANY"`))
	}

	bearerTokenPattern := regexp.MustCompile(`"Authorization":\s*\[\s*"\s*Bearer\s+[A-Za-z0-9-_\.]+"\s*\]`)
	jsonData = bearerTokenPattern.ReplaceAll(jsonData, []byte(`"Authorization": ["Bearer ACCESS_TOKEN"]`))

	return jsonData
}
