package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"testing"
	"time"
)

func TestRepository_AddAction(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection: %s", err, mock)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO user_activities (id, action_id, product_id, created_at) VALUES ($1, $2, $3, $4)`).
		WithArgs(sqlmock.AnyArg(), "00000000-0000-1000-0000-000000000000", "00000000-0000-0000-0000-000000000001", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repository := &Repository{db: db}

	err = repository.AddAction("00000000-0000-1000-0000-000000000000", "00000000-0000-0000-0000-000000000001")
	if err != nil {
		t.Errorf("should add action to the db, receive: %T\n", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestRepository_ActionWithInterval(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection: %s", err, mock)
	}
	defer db.Close()

	date, _ := time.Parse("2022-08-31T14:29:06.484648Z", "2022-08-31T14:29:06.484648Z")

	mock.ExpectQuery(`SELECT products.*, user_activities.action_id, user_activities.created_at FROM products JOIN user_activities on (
    	user_activities.action_id=$1 AND 
    	user_activities.product_id=products.id) WHERE 
    	user_activities.created_at between $2 and $3;`).
		WithArgs("00000000-0000-1000-0000-000000000000", "2022-08-30", "2022-09-01").
		WillReturnRows(
			mock.
				NewRows([]string{"productID", "name", "description", "price", "category", "actionID", "createdAt"}).
				AddRow(
					"00000000-0000-0000-0000-000000000001",
					"Shampoo",
					"Gel",
					100,
					"00000000-0000-0000-1000-000000000000",
					"00000000-0000-1000-0000-000000000000",
					date),
		)

	repositoriy := &Repository{db: db}

	res, err := repositoriy.ActionWithInterval("00000000-0000-1000-0000-000000000000", "2022-08-30", "2022-09-01")
	if err != nil {
		t.Errorf("unexpected db err: %s\n", err)
		return
	}

	exp := []model.DBResponse{
		{
			"00000000-0000-1000-0000-000000000000",
			date,
			"00000000-0000-0000-0000-000000000001",
			"Shampoo",
			"Gel",
			100,
			"00000000-0000-0000-1000-000000000000"},
	}

	if !testDBResponseEqual(res, exp) {
		t.Errorf("should be equal, but got: %v, want: %v\n", res, exp)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func testDBResponseEqual(a, b []model.DBResponse) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRepository_TwoActionsWithInterval(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection: %s", err, mock)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT products.* FROM products JOIN user_activities as ua1 on (
    	ua1.action_id=$1 AND ua1.product_id=products.id)
    	JOIN user_activities as ua2 on (
    	ua2.action_id=$2 AND ua2.product_id=products.id) WHERE 
    	ua1.created_at between $3 and $4;`).
		WithArgs("00000000-0000-1000-0000-000000000000", "00000000-0000-2000-0000-000000000000", "2022-08-30", "2022-09-01").
		WillReturnRows(
			mock.
				NewRows([]string{"productID", "name", "description", "price", "category"}).
				AddRow(
					"00000000-0000-0000-0000-000000000003",
					"Toothpaste",
					"Colgate",
					15.5,
					"00000000-0000-0000-1000-000000000000"),
		)

	repositoriy := &Repository{db: db}

	res, err := repositoriy.TwoActionsWithInterval("00000000-0000-1000-0000-000000000000", "00000000-0000-2000-0000-000000000000", "2022-08-30", "2022-09-01")
	if err != nil {
		t.Errorf("unexpected db err: %s\n", err)
		return
	}

	exp := []model.DBResponse2Actions{
		{
			"00000000-0000-0000-0000-000000000003",
			"Toothpaste",
			"Colgate",
			15.5,
			"00000000-0000-0000-1000-000000000000"},
	}

	if !testDB2ActionsEqual(res, exp) {
		t.Errorf("should be equal, but got: %v, want: %v\n", res, exp)
	}
}

func testDB2ActionsEqual(a, b []model.DBResponse2Actions) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
