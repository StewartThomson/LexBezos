package main

import (
	"reflect"
	"testing"
)

func Test_cleanURL(t *testing.T) {
	type args struct {
		inputUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "some qs params",
			args: args{
				inputUrl: "https://cnn.com/blah-blah?k=v&foo=bar#",
			},
			want:    "https://cnn.com/blah-blah",
			wantErr: false,
		},
		{

			name: "no qs params",
			args: args{
				inputUrl: "https://cnn.com/blah-blah",
			},
			want:    "https://cnn.com/blah-blah",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cleanURL(tt.args.inputUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("cleanURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cleanURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_modifyTitles(t *testing.T) {
	type args struct {
		listings []Listing
	}
	tests := []struct {
		name         string
		args         args
		wantModified []Listing
	}{
		{
			name: "Buncha replacements",
			args: args{
				listings: []Listing{
					Listing{
						ListingData{
							Title: "Jeff Bezos' Amazon amazon event",
						},
					},
				},
			},
			wantModified: []Listing{
				Listing{
					ListingData{
						Title: "Lex Luthor's LexCorp LexCorp event",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotModified := modifyTitles(tt.args.listings); !reflect.DeepEqual(gotModified, tt.wantModified) {
				t.Errorf("modifyTitles() = %v, want %v", gotModified, tt.wantModified)
			}
		})
	}
}
