package link

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		r io.Reader
	}
	f, err := os.Open("ex1.html")
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}
	defer f.Close()

	tests := []struct {
		name    string
		args    args
		want    []Link
		wantErr bool
	}{
		{
			name: "Parse ex1.html",
			args: args{r: f},
			want: []Link{
				{
					Href: "/other-page",
					Text: "A link to another page",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
