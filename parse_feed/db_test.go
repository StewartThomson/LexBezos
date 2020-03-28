package main

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

func TestPopulateTweetTable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectPrepare("INSERT INTO LexBezos.tweets")
	mock.ExpectExec("INSERT INTO LexBezos.tweets").WithArgs(69, "hello world").WillReturnResult(sqlmock.NewResult(1, 1))

	type args struct {
		db       *sql.DB
		listings []Listing
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ensure they're going in right",
			args: args{
				db: db,
				listings: []Listing{
					{
						ListingData{
							Title: "hello",
							Url:   "world",
							DBID:  69,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PopulateTweetTable(tt.args.db, tt.args.listings); (err != nil) != tt.wantErr {
				t.Errorf("PopulateTweetTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_filterPostedListings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectPrepare("SELECT id FROM LexBezos.articles")
	mock.ExpectQuery("SELECT id FROM LexBezos.articles").WithArgs("world").WillReturnRows(mock.NewRows([]string{"id"}).FromCSVString("1"))
	mock.ExpectQuery("SELECT id FROM LexBezos.articles").WithArgs("world1").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT id FROM LexBezos.articles").WithArgs("world2").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT title FROM LexBezos.articles").WillReturnRows(mock.NewRows([]string{"title"}).FromCSVString("hello"))
	mock.ExpectQuery("SELECT title FROM LexBezos.articles").WillReturnRows(mock.NewRows([]string{"title"}).FromCSVString("hopefully nothing alike"))
	type args struct {
		db       *sql.DB
		listings []Listing
	}
	tests := []struct {
		name                 string
		args                 args
		wantApprovedListings []Listing
		wantBadTitleListings []Listing
		wantErr              bool
	}{
		{
			name: "random post",
			args: args{
				db: db,
				listings: []Listing{
					{
						ListingData{
							Title: "hello",
							Url:   "world",
							DBID:  0,
						},
					},
					{
						ListingData{
							Title: "hello1",
							Url:   "world1",
							DBID:  0,
						},
					},
					{
						ListingData{
							Title: "something complete different",
							Url:   "world2",
							DBID:  0,
						},
					},
				},
			},
			wantApprovedListings: []Listing{
				{
					ListingData{
						Title: "something complete different",
						Url:   "world2",
						DBID:  0,
					},
				},
			},
			wantBadTitleListings: []Listing{
				{
					ListingData{
						Title: "hello1",
						Url:   "world1",
						DBID:  0,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotApprovedListings, gotBadTitleListings, err := filterPostedListings(tt.args.db, tt.args.listings)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterPostedListings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotApprovedListings, tt.wantApprovedListings) {
				t.Errorf("filterPostedListings() gotApprovedListings = %v, want %v", gotApprovedListings, tt.wantApprovedListings)
			}
			if !reflect.DeepEqual(gotBadTitleListings, tt.wantBadTitleListings) {
				t.Errorf("filterPostedListings() gotBadTitleListings = %v, want %v", gotBadTitleListings, tt.wantBadTitleListings)
			}
		})
	}
}

func Test_storeListings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectPrepare("INSERT INTO LexBezos.articles")
	mock.ExpectExec("INSERT INTO LexBezos.articles").WithArgs("world", "hello").WillReturnResult(sqlmock.NewResult(1, 1))

	type args struct {
		db       *sql.DB
		listings []Listing
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{
			name: "Ensure they're going in right",
			args: args{
				db: db,
				listings: []Listing{
					{
						ListingData{
							Title: "hello",
							Url:   "world",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storeListings(tt.args.db, tt.args.listings); (err != nil) != tt.wantErr {
				t.Errorf("storeListings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
