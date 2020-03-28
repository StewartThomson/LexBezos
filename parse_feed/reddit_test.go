package main

import (
	"reflect"
	"testing"
)

func Test_getJeffListings(t *testing.T) {
	type args struct {
		listings []Listing
	}
	tests := []struct {
		name string
		args args
		want []Listing
	}{
		{
			name: "A few posts",
			args: args{[]Listing{
				Listing{
					ListingData{
						Title: "Bezos disgustingly rich",
					},
				},
				Listing{
					ListingData{
						Title: "Gates disgustingly rich",
					},
				},
				Listing{
					ListingData{
						Title: "Amazon burns",
					},
				},
				Listing{
					ListingData{
						Title: "Amazon soars",
					},
				},
			},
			},
			want: []Listing{
				Listing{
					ListingData{
						Title: "Bezos disgustingly rich",
					},
				},
				Listing{
					ListingData{
						Title: "Amazon soars",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getJeffListings(tt.args.listings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getJeffListings() = %v, want %v", got, tt.want)
			}
		})
	}
}
