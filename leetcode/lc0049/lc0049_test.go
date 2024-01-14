package lc0049

import (
	"log"
	"reflect"
	"sort"
	"testing"
)

type args struct {
	strs []string
}

var tests = []struct {
	name string
	args args
	want [][]string
}{
	{
		name: "case1",
		args: args{
			strs: []string{"eat", "tea", "tan", "ate", "nat", "bat"},
		},
		want: [][]string{
			{"eat", "tea", "ate"},
			{"tan", "nat"},
			{"bat"},
		},
	},
	{
		name: "case2",
		args: args{
			strs: []string{""},
		},
		want: [][]string{
			{""},
		},
	},
	{
		name: "case3",
		args: args{
			strs: []string{"a"},
		},
		want: [][]string{
			{"a"},
		},
	},
	{
		name: "case4",
		args: args{
			strs: []string{"ddddddddddg", "dgggggggggg"},
		},
		want: [][]string{
			{"ddddddddddg"},
			{"dgggggggggg"},
		},
	},
}

func Test_groupAnagrams(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupAnagrams(tt.args.strs)
			sortSlice(got)
			sortSlice(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupAnagrams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupAnagramsPro(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupAnagramsPro(tt.args.strs)
			sortSlice(got)
			sortSlice(tt.want)
			//切片元素顺序不同，DeepEqual也会返回不同
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupAnagrams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sortSlice(s [][]string) {
	for i := range s {
		sort.Strings(s[i])
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i][0] <= s[j][0]
	})
	log.Println("s sorted", s)
}
