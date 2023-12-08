package almanac

import (
	"reflect"
	"testing"
)

func TestSeedMap_LocationForRange(t *testing.T) {
	type fields struct {
		From         string
		To           string
		Ranges       [][2]int
		Destinations []int
	}
	type args struct {
		ranges [][2]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][2]int
	}{
		{
			name: "end",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}},
				Destinations: []int{50},
			},
			args: args{
				ranges: [][2]int{{10, 20}},
			},
			want: [][2]int{{20, 20}, {50, 59}},
		},
		{
			name: "end range",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}},
				Destinations: []int{50},
			},
			args: args{
				ranges: [][2]int{{10, 21}},
			},
			want: [][2]int{{20, 21}, {50, 59}},
		},
		{
			name: "end exact",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}},
				Destinations: []int{50},
			},
			args: args{
				ranges: [][2]int{{15, 19}},
			},
			want: [][2]int{{55, 59}},
		},
		{
			name: "start exact",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}},
				Destinations: []int{50},
			},
			args: args{
				ranges: [][2]int{{10, 15}},
			},
			want: [][2]int{{50, 55}},
		},
		{
			name: "start",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}},
				Destinations: []int{50},
			},
			args: args{
				ranges: [][2]int{{9, 11}},
			},
			want: [][2]int{{9, 9}, {50, 51}},
		},
		{
			name: "start range",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}},
				Destinations: []int{20},
			},
			args: args{
				ranges: [][2]int{{18, 20}},
			},
			want: [][2]int{{20, 20}, {28, 29}},
		},
		{
			name: "mid",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}},
				Destinations: []int{50},
			},
			args: args{
				ranges: [][2]int{{8, 21}},
			},
			want: [][2]int{{8, 9}, {20, 21}, {50, 59}},
		},
		{
			name: "miss",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}},
				Destinations: []int{50},
			},
			args: args{
				ranges: [][2]int{{29, 39}},
			},
			want: [][2]int{{29, 39}},
		},
		{
			name: "multiple",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}, {20, 29}},
				Destinations: []int{50, 100},
			},
			args: args{
				ranges: [][2]int{{19, 29}},
			},
			want: [][2]int{{59, 59}, {100, 109}},
		},
		{
			name: "multiple args",
			fields: fields{
				From:         "seed",
				To:           "soil",
				Ranges:       [][2]int{{10, 19}, {20, 29}},
				Destinations: []int{50, 100},
			},
			args: args{
				ranges: [][2]int{{19, 29}, {31, 39}},
			},
			want: [][2]int{{20, 29}, {59, 59}, {19, 19}, {100, 109}, {31, 39}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := SeedMap{
				From:         tt.fields.From,
				To:           tt.fields.To,
				Ranges:       tt.fields.Ranges,
				Destinations: tt.fields.Destinations,
			}
			if got := sm.LocationForRange(tt.args.ranges); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SeedMap.LocationForRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
