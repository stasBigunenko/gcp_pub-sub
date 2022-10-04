package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"Intern/gcp_pub-sub/modules/subscriber/service"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var _ service.Service = mockService{}

type mockService struct {
	mockActionIDWithInterval     func(date model.InputWithDate) ([]model.DBResponse, error)
	mockTwoActionsIDWithInterval func(model.InputWithDate) ([]model.DBResponse2Actions, error)
}

func (m mockService) ActionIDWithInterval(date model.InputWithDate) ([]model.DBResponse, error) {
	return m.mockActionIDWithInterval(date)
}

func (m mockService) TwoActionsIDWithInterval(date model.InputWithDate) ([]model.DBResponse2Actions, error) {
	return m.mockTwoActionsIDWithInterval(date)
}

func TestProduct_ProductsInBucket(t *testing.T) {
	date, _ := time.Parse("2022-08-31T14:29:06.484648Z", "2022-08-31T14:29:06.484648Z")

	type fields struct {
		s service.Service
	}
	type args struct {
		body   []byte
		status int
		want   []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "correct test",
			fields: fields{
				s: mockService{
					mockActionIDWithInterval: func(model.InputWithDate) ([]model.DBResponse, error) {
						return []model.DBResponse{
							{"00000000-0000-1000-0000-000000000000",
								date,
								"00000000-0000-0000-0000-000000000001",
								"Shampoo",
								"Gel",
								100,
								"00000000-0000-0000-1000-000000000000"},
						}, nil
					},
				},
			},
			args: args{
				body: []byte(`{
						"actionID":"00000000-0000-1000-0000-000000000000",
    					"fromDateYear":"2022",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30"
						}`),
				want: []byte(
					`{"products":[{"actionID":"00000000-0000-1000-0000-000000000000","createdAt":"0001-01-01T00:00:00Z","productID":"00000000-0000-0000-0000-000000000001","name":"Shampoo","description":"Gel","price":100,"category":"00000000-0000-0000-1000-000000000000"}]}`,
				),
				status: http.StatusOK,
			},
		},
		{
			name: "service error test",
			fields: fields{
				s: mockService{
					mockActionIDWithInterval: func(model.InputWithDate) ([]model.DBResponse, error) {
						return []model.DBResponse{}, errors.New("service error")
					},
				},
			},
			args: args{
				body: []byte(`{
						"actionID":"00000000-0000-1000-0000-000000000000",
    					"fromDateYear":"2022",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30"
						}`),
				want: []byte(
					`{"error":"service error"}`,
				),
				status: http.StatusBadRequest,
			},
		},
		{
			name: "json error test",
			fields: fields{
				s: mockService{
					mockActionIDWithInterval: func(model.InputWithDate) ([]model.DBResponse, error) {
						return []model.DBResponse{}, nil
					},
				},
			},
			args: args{
				body: []byte(`{
						"action":"00000000-0000-1000-0000-000000000000",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30""
						}`),
				want: []byte(
					`{"error":"json error"}`,
				),
				status: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			p := &Product{
				service: tt.fields.s,
			}

			router := gin.Default()

			router.Use(func(c *gin.Context) {
				c.Set("input", tt.args.body)
			}).GET("/bucket", p.ProductsInBucket)

			rr := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/bucket", bytes.NewBuffer(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}

			router.ServeHTTP(rr, req)

			if tt.args.status != rr.Code {
				t.Errorf("error with Status: want: %v, got %v\n", tt.args.status, rr.Code)
				return
			}

			if !checkEqualityDBResponse(tt.args.want, rr.Body.Bytes()) {
				t.Errorf("error with body: want %s, got %s\n", tt.args.want, rr.Body.Bytes())
				return
			}
		})
	}
}

func TestProduct_ProductsOutFromBucket(t *testing.T) {
	date, _ := time.Parse("2022-08-31T14:29:06.484648Z", "2022-08-31T14:29:06.484648Z")

	type fields struct {
		s service.Service
	}
	type args struct {
		body   []byte
		status int
		want   []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "correct test",
			fields: fields{
				s: mockService{
					mockActionIDWithInterval: func(model.InputWithDate) ([]model.DBResponse, error) {
						return []model.DBResponse{
							{"00000000-0000-1000-0000-000000000000",
								date,
								"00000000-0000-0000-0000-000000000001",
								"Shampoo",
								"Gel",
								100,
								"00000000-0000-0000-1000-000000000000"},
						}, nil
					},
				},
			},
			args: args{
				body: []byte(`{
						"actionID":"00000000-0000-1000-0000-000000000000",
    					"fromDateYear":"2022",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30"
						}`),
				want: []byte(
					`{"products":[{"actionID":"00000000-0000-1000-0000-000000000000","createdAt":"0001-01-01T00:00:00Z","productID":"00000000-0000-0000-0000-000000000001","name":"Shampoo","description":"Gel","price":100,"category":"00000000-0000-0000-1000-000000000000"}]}`,
				),
				status: http.StatusOK,
			},
		},
		{
			name: "service error test",
			fields: fields{
				s: mockService{
					mockActionIDWithInterval: func(model.InputWithDate) ([]model.DBResponse, error) {
						return []model.DBResponse{}, errors.New("service error")
					},
				},
			},
			args: args{
				body: []byte(`{
						"actionID":"00000000-0000-1000-0000-000000000000",
    					"fromDateYear":"2022",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30"
						}`),
				want: []byte(
					`{"error":"service error"}`,
				),
				status: http.StatusBadRequest,
			},
		},
		{
			name: "json error test",
			fields: fields{
				s: mockService{
					mockActionIDWithInterval: func(model.InputWithDate) ([]model.DBResponse, error) {
						return []model.DBResponse{}, nil
					},
				},
			},
			args: args{
				body: []byte(`{
						"action":"00000000-0000-1000-0000-000000000000",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30""
						}`),
				want: []byte(
					`{"error":"json error"}`,
				),
				status: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			p := &Product{
				service: tt.fields.s,
			}

			router := gin.Default()

			router.Use(func(c *gin.Context) {
				c.Set("input", tt.args.body)
			}).GET("/outofbucket", p.ProductsOutFromBucket)

			rr := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/outofbucket", bytes.NewBuffer(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}

			router.ServeHTTP(rr, req)

			if tt.args.status != rr.Code {
				t.Errorf("error with Status: %v\n", err)
				return
			}

			if !checkEqualityDBResponse(tt.args.want, rr.Body.Bytes()) {
				t.Errorf("error with body: want %s, got %s\n", tt.args.want, rr.Body.Bytes())
				return
			}
		})
	}
}

func TestProduct_ProductsDescription(t *testing.T) {
	date, _ := time.Parse("2022-08-31T14:29:06.484648Z", "2022-08-31T14:29:06.484648Z")

	type fields struct {
		s service.Service
	}
	type args struct {
		body   []byte
		status int
		want   []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "correct test",
			fields: fields{
				s: mockService{
					mockActionIDWithInterval: func(model.InputWithDate) ([]model.DBResponse, error) {
						return []model.DBResponse{
							{"00000000-0000-1000-0000-000000000000",
								date,
								"00000000-0000-0000-0000-000000000001",
								"Shampoo",
								"Gel",
								100,
								"00000000-0000-0000-1000-000000000000"},
						}, nil
					},
				},
			},
			args: args{
				body: []byte(`{
						"actionID":"00000000-0000-1000-0000-000000000000",
    					"fromDateYear":"2022",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30"
						}`),
				want: []byte(
					`{"products":[{"actionID":"00000000-0000-1000-0000-000000000000","createdAt":"0001-01-01T00:00:00Z","productID":"00000000-0000-0000-0000-000000000001","name":"Shampoo","description":"Gel","price":100,"category":"00000000-0000-0000-1000-000000000000"}]}`,
				),
				status: http.StatusOK,
			},
		},
		{
			name: "service error test",
			fields: fields{
				s: mockService{
					mockActionIDWithInterval: func(model.InputWithDate) ([]model.DBResponse, error) {
						return []model.DBResponse{}, errors.New("service error")
					},
				},
			},
			args: args{
				body: []byte(`{
						"actionID":"00000000-0000-1000-0000-000000000000",
    					"fromDateYear":"2022",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30"
						}`),
				want: []byte(
					`{"error":"service error"}`,
				),
				status: http.StatusBadRequest,
			},
		},
		{
			name: "json error test",
			fields: fields{
				s: mockService{
					mockActionIDWithInterval: func(model.InputWithDate) ([]model.DBResponse, error) {
						return []model.DBResponse{}, nil
					},
				},
			},
			args: args{
				body: []byte(`{
						"action":"00000000-0000-1000-0000-000000000000",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30""
						}`),
				want: []byte(
					`{"error":"json error"}`,
				),
				status: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			p := &Product{
				service: tt.fields.s,
			}

			router := gin.Default()

			router.Use(func(c *gin.Context) {
				c.Set("input", tt.args.body)
			}).GET("/description", p.ProductsDescription)

			rr := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/description", bytes.NewBuffer(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}

			router.ServeHTTP(rr, req)

			if tt.args.status != rr.Code {
				t.Errorf("error with Status: %v\n", err)
				return
			}

			if !checkEqualityDBResponse(tt.args.want, rr.Body.Bytes()) {
				t.Errorf("error with body: want %s, got %s\n", tt.args.want, rr.Body.Bytes())
				return
			}
		})
	}
}

func TestProduct_ProductsBucketAndDescription(t *testing.T) {
	type fields struct {
		s service.Service
	}
	type args struct {
		body   []byte
		status int
		want   []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "correct test",
			fields: fields{
				s: mockService{
					mockTwoActionsIDWithInterval: func(model.InputWithDate) ([]model.DBResponse2Actions, error) {
						return []model.DBResponse2Actions{
							{"00000000-0000-0000-0000-000000000001",
								"Shampoo",
								"Gel",
								100,
								"00000000-0000-0000-1000-000000000000"},
						}, nil
					},
				},
			},
			args: args{
				body: []byte(`{
						"actionID":"00000000-0000-1000-0000-000000000000",
						"actionID2":"00000000-0000-2000-0000-000000000000",
    					"fromDateYear":"2022",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30"
						}`),
				want: []byte(
					`{"products":[{"productID":"00000000-0000-0000-0000-000000000001","name":"Shampoo","description":"Gel","price":100,"category":"00000000-0000-0000-1000-000000000000"}]}`,
				),
				status: http.StatusOK,
			},
		},
		{
			name: "service error test",
			fields: fields{
				s: mockService{
					mockTwoActionsIDWithInterval: func(model.InputWithDate) ([]model.DBResponse2Actions, error) {
						return []model.DBResponse2Actions{}, errors.New("service error")
					},
				},
			},
			args: args{
				body: []byte(`{
						"actionID":"00000000-0000-1000-0000-000000000000",
						"actionID2":"00000000-0000-2000-0000-000000000000",
    					"fromDateYear":"2022",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30"
						}`),
				want: []byte(
					`{"error":"service error"}`,
				),
				status: http.StatusBadRequest,
			},
		},
		{
			name: "json error test",
			fields: fields{
				s: mockService{
					mockTwoActionsIDWithInterval: func(model.InputWithDate) ([]model.DBResponse2Actions, error) {
						return []model.DBResponse2Actions{}, nil
					},
				},
			},
			args: args{
				body: []byte(`{
						"actionID":"00000000-0000-1000-0000-000000000000",
						"actionID2":"00000000-0000-2000-0000-000000000000"
    					"fromDateYear":"2022",
    					"fromDateMonth":"08",
    					"fromDateDay":"24",
    					"toDateYear":"2022",
    					"toDateMonth":"09",
    					"toDateDay":"30""
						}`),
				want: []byte(
					`{"error":"json error"}`,
				),
				status: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			p := &Product{
				service: tt.fields.s,
			}

			router := gin.Default()

			router.Use(func(c *gin.Context) {
				c.Set("input", tt.args.body)
			}).GET("/descriptionandbucket", p.ProductsBucketAndDescription)

			rr := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/descriptionandbucket", bytes.NewBuffer(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}

			router.ServeHTTP(rr, req)

			if tt.args.status != rr.Code {
				t.Errorf("error with Status: %v\n", err)
				return
			}

			if !checkEqualityDBResponse2Action(tt.args.want, rr.Body.Bytes()) {
				t.Errorf("error with body: want %s, got %s\n", tt.args.want, rr.Body.Bytes())
				return
			}
		})
	}
}

func checkEqualityDBResponse(want, res []byte) bool {
	var (
		wantModel map[string][]model.DBResponse
		resModel  map[string][]model.DBResponse
	)

	json.Unmarshal(want, &wantModel)
	json.Unmarshal(res, &resModel)

	for i := range wantModel {
		for j := 0; j < len(wantModel[i]); j++ {
			if wantModel[i][j] != resModel[i][j] {
				return false
			}
		}
	}

	return true
}

func checkEqualityDBResponse2Action(want, res []byte) bool {
	var (
		wantModel map[string][]model.DBResponse2Actions
		resModel  map[string][]model.DBResponse2Actions
	)

	json.Unmarshal(want, &wantModel)
	json.Unmarshal(res, &resModel)

	for i := range wantModel {
		for j := 0; j < len(wantModel[i]); j++ {
			if wantModel[i][j] != resModel[i][j] {
				return false
			}
		}
	}

	return true
}
