package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"Intern/gcp_pub-sub/modules/subscriber/repo"
	"errors"
	"testing"
	"time"
)

var _ repo.ProductsRepo = mockRepo{}

type mockRepo struct {
	mockAddAction func(string, string) error

	mockActionWithInterval     func(string, string, string) ([]model.DBResponse, error)
	mockTwoActionsWithInterval func(string, string, string, string) ([]model.DBResponse2Actions, error)
}

func (m mockRepo) AddAction(actionId, productId string) error {
	return m.mockAddAction(actionId, productId)
}

func (m mockRepo) ActionWithInterval(actionID, fromDate, toDate string) ([]model.DBResponse, error) {
	return m.mockActionWithInterval(actionID, fromDate, toDate)
}

func (m mockRepo) TwoActionsWithInterval(actionID, actionID2, fromDate, toDate string) ([]model.DBResponse2Actions, error) {
	return m.mockTwoActionsWithInterval(actionID, actionID2, fromDate, toDate)
}

func checkEqualityDBResponse(res, want []model.DBResponse) bool {
	for i := 0; i < len(res); i++ {
		if res[i] != want[i] {
			return false
		}
	}
	return true
}

func checkEqualityDBResponse2Action(res, want []model.DBResponse2Actions) bool {
	for i := 0; i < len(res); i++ {
		if res[i] != want[i] {
			return false
		}
	}
	return true
}

func TestServiceProduct_ActionIDWithInterval(t *testing.T) {
	date, _ := time.Parse("2022-08-31T14:29:06.484648Z", "2022-08-31T14:29:06.484648Z")

	type fields struct {
		r repo.ProductsRepo
	}
	type args struct {
		input   model.InputWithDate
		want    []model.DBResponse
		wantErr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Correct test",
			fields: fields{
				r: mockRepo{
					mockActionWithInterval: func(actionID, fromDate, toDate string) ([]model.DBResponse, error) {
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
				input: model.InputWithDate{
					ActionID:      "00000000-0000-1000-0000-000000000000",
					FromDateYear:  "2022",
					FromDateMonth: "08",
					FromDateDay:   "24",
					ToDateYear:    "2022",
					ToDateMonth:   "09",
					ToDateDay:     "30"},
				want: []model.DBResponse{
					{"00000000-0000-1000-0000-000000000000",
						date,
						"00000000-0000-0000-0000-000000000001",
						"Shampoo",
						"Gel",
						100,
						"00000000-0000-0000-1000-000000000000"},
				},
				wantErr: "",
			},
		},
		{
			name: "Error db test",
			fields: fields{
				r: mockRepo{
					mockActionWithInterval: func(actionID, fromDate, toDate string) ([]model.DBResponse, error) {
						return []model.DBResponse{}, errors.New("db problem")
					},
				},
			},
			args: args{
				input: model.InputWithDate{
					ActionID:      "00000000-0000-1000-0000-000000000000",
					FromDateYear:  "2022",
					FromDateMonth: "08",
					FromDateDay:   "24",
					ToDateYear:    "2022",
					ToDateMonth:   "08",
					ToDateDay:     "25"},
				want:    []model.DBResponse{},
				wantErr: "db problem",
			},
		},
		{
			name: "Error validate date test",
			fields: fields{
				r: mockRepo{
					mockActionWithInterval: func(actionID, fromDate, toDate string) ([]model.DBResponse, error) {
						return []model.DBResponse{}, errors.New("db problem")
					},
				},
			},
			args: args{
				input: model.InputWithDate{
					ActionID:      "00000000-0000-1000-0000-000000000000",
					FromDateYear:  "202",
					FromDateMonth: "13",
					FromDateDay:   "41",
					ToDateYear:    "22",
					ToDateMonth:   "08",
					ToDateDay:     "25"},
				want:    []model.DBResponse{},
				wantErr: `parsing time "202-13-41" as "2006-01-02": cannot parse "13-41" as "2006"`,
			},
		},
		{
			name: "Error validate date test2",
			fields: fields{
				r: mockRepo{
					mockActionWithInterval: func(actionID, fromDate, toDate string) ([]model.DBResponse, error) {
						return []model.DBResponse{}, errors.New("db problem")
					},
				},
			},
			args: args{
				input: model.InputWithDate{
					ActionID:      "00000000-0000-1000-0000-000000000000",
					FromDateYear:  "2020",
					FromDateMonth: "08",
					FromDateDay:   "24",
					ToDateYear:    "22",
					ToDateMonth:   "08",
					ToDateDay:     "25"},
				want:    []model.DBResponse{},
				wantErr: `parsing time "22-08-25" as "2006-01-02": cannot parse "8-25" as "2006"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ServiceProduct{
				repo: tt.fields.r,
			}

			res, err := s.ActionIDWithInterval(tt.args.input)
			if err != nil {
				if err.Error() != tt.args.wantErr {
					t.Errorf("unexpected error: want: %s, got %s\n", tt.args.wantErr, err)
					return
				}
			}

			if !checkEqualityDBResponse(res, tt.args.want) && len(res) > 0 {
				t.Errorf("res and want are not equal: want: %v, got %v\n", tt.args.want, res)
				return
			}
		})
	}
}

func TestServiceProduct_TwoActionsIDWithInterval(t *testing.T) {
	type fields struct {
		r repo.ProductsRepo
	}
	type args struct {
		input   model.InputWithDate
		want    []model.DBResponse2Actions
		wantErr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Correct test",
			fields: fields{
				r: mockRepo{
					mockTwoActionsWithInterval: func(actionID, actionID2, fromDate, toDate string) ([]model.DBResponse2Actions, error) {
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
				input: model.InputWithDate{
					ActionID:      "00000000-0000-1000-0000-000000000000",
					ActionID2:     "00000000-0000-2000-0000-000000000000",
					FromDateYear:  "2022",
					FromDateMonth: "08",
					FromDateDay:   "24",
					ToDateYear:    "2022",
					ToDateMonth:   "09",
					ToDateDay:     "30"},
				want: []model.DBResponse2Actions{
					{"00000000-0000-0000-0000-000000000001",
						"Shampoo",
						"Gel",
						100,
						"00000000-0000-0000-1000-000000000000"},
				},
				wantErr: "",
			},
		},
		{
			name: "Error db test",
			fields: fields{
				r: mockRepo{
					mockTwoActionsWithInterval: func(actionID, actionID2, fromDate, toDate string) ([]model.DBResponse2Actions, error) {
						return []model.DBResponse2Actions{}, errors.New("db problem")
					},
				},
			},
			args: args{
				input: model.InputWithDate{
					ActionID:      "00000000-0000-1000-0000-000000000000",
					ActionID2:     "00000000-0000-2000-0000-000000000000",
					FromDateYear:  "2022",
					FromDateMonth: "08",
					FromDateDay:   "24",
					ToDateYear:    "2022",
					ToDateMonth:   "08",
					ToDateDay:     "25"},
				want:    []model.DBResponse2Actions{},
				wantErr: "db problem",
			},
		},
		{
			name: "Error validate date test",
			fields: fields{
				r: mockRepo{
					mockTwoActionsWithInterval: func(actionID, actionID2, fromDate, toDate string) ([]model.DBResponse2Actions, error) {
						return []model.DBResponse2Actions{}, errors.New("db problem")
					},
				},
			},
			args: args{
				input: model.InputWithDate{
					ActionID:      "00000000-0000-1000-0000-000000000000",
					FromDateYear:  "202",
					FromDateMonth: "13",
					FromDateDay:   "41",
					ToDateYear:    "22",
					ToDateMonth:   "08",
					ToDateDay:     "25"},
				want:    []model.DBResponse2Actions{},
				wantErr: `parsing time "202-13-41" as "2006-01-02": cannot parse "13-41" as "2006"`,
			},
		},
		{
			name: "Error validate date test2",
			fields: fields{
				r: mockRepo{
					mockTwoActionsWithInterval: func(actionID, actionID2, fromDate, toDate string) ([]model.DBResponse2Actions, error) {
						return []model.DBResponse2Actions{}, errors.New("db problem")
					},
				},
			},
			args: args{
				input: model.InputWithDate{
					ActionID:      "00000000-0000-1000-0000-000000000000",
					FromDateYear:  "2020",
					FromDateMonth: "08",
					FromDateDay:   "24",
					ToDateYear:    "22",
					ToDateMonth:   "08",
					ToDateDay:     "25"},
				want:    []model.DBResponse2Actions{},
				wantErr: `parsing time "22-08-25" as "2006-01-02": cannot parse "8-25" as "2006"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ServiceProduct{
				repo: tt.fields.r,
			}

			res, err := s.TwoActionsIDWithInterval(tt.args.input)
			if err != nil {
				if err.Error() != tt.args.wantErr {
					t.Errorf("unexpected error: want: %s, got %s\n", tt.args.wantErr, err)
					return
				}
			}

			if !checkEqualityDBResponse2Action(res, tt.args.want) && len(res) > 0 {
				t.Errorf("res and want are not equal: want: %v, got %v\n", tt.args.want, res)
				return
			}
		})
	}
}
