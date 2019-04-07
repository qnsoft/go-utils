package maputil

import (
	"fmt"
	"testing"
	"time"
)

type Person struct {
	Id         int               `json:"id"`
	Name       string            `json:"name"`
	Age        int               `json:"age"`
	Sex        bool              `json:"sex"`
	LikeMap    map[string]string `json:"like_map"`
	CreateTime time.Time         `json:"create_time"`
}

type Student struct {
	Person      Person   `json:"person"`
	Class       string   `json:"class"`
	CourseArray []string `json:"course_array"`
}

var p = Person{
	Id:         123456,
	Name:       "zhang san",
	Age:        23,
	Sex:        true,
	LikeMap:    map[string]string{"花": "yes", "宠物": "yes"},
	CreateTime: time.Now(),
}

var student = Student{
	Person:      p,
	Class:       "二班",
	CourseArray: []string{"语", "数", "化"},
}

func TestStructToMap(t *testing.T) {
	type args struct {
		a interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				a: student,
			},
		}, {
			name: "ptr",
			args: args{
				a: &student,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StructToMapReflect(tt.args.a)
			fmt.Println(got)
		})
	}
}

func TestStructToMapJson(t *testing.T) {
	type args struct {
		o interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				o: student,
			},
		}, {
			name: "ptr",
			args: args{
				o: &student,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StructToMapJson(tt.args.o)
			fmt.Println(got)
		})
	}
}
