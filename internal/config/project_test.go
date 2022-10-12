package config

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"testing"
)

func TestProject_GetStates(t *testing.T) {
	type fields struct {
		Terraform map[string]*Terraform
	}
	type args struct {
		stacks []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*interface{}
	}{
		{name: "success", fields: fields{
			Terraform: map[string]*Terraform{
				"test":  {},
				"test1": {},
				"test2": {},
				"test3": {},
			},
		}, args: args{}, want: map[string]*interface{}{"test": nil, "test1": nil, "test2": nil, "test3": nil}},
		{name: "success", fields: fields{
			Terraform: map[string]*Terraform{
				"test":  {},
				"test1": {},
				"test2": {},
				"test3": {},
			},
		}, args: args{stacks: []string{"test", "test2"}}, want: map[string]*interface{}{"test": nil, "test2": nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				Terraform: tt.fields.Terraform,
			}
			got := p.GetStates(tt.args.stacks...)
			gotKeys := maps.Keys(got)
			wantKeys := maps.Keys(tt.want)
			slices.Sort(gotKeys)
			slices.Sort(wantKeys)
			if !slices.Equal(gotKeys, wantKeys) {
				t.Errorf("GetStates() = %v, want %v", maps.Keys(got), maps.Keys(tt.want))
			}
		})
	}
}

func TestProject_GetApps(t *testing.T) {
	type fields struct {
		Serverless map[string]*Serverless
		Ecs        map[string]*Ecs
		Alias      map[string]*Alias
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*interface{}
	}{
		{name: "success", fields: fields{
			Serverless: map[string]*Serverless{
				"testSls1": {},
				"testSls2": {},
			},
			Ecs: map[string]*Ecs{
				"testEcs1": {},
				"testEcs2": {},
			},
			Alias: map[string]*Alias{
				"testAlias1": {},
				"testAlias2": {},
			},
		}, args: args{}, want: map[string]*interface{}{"testAlias1": nil, "testAlias2": nil, "testEcs1": nil, "testEcs2": nil, "testSls2": nil, "testSls1": nil}},
		{name: "success", fields: fields{
			Serverless: map[string]*Serverless{
				"testSls1": {},
				"testSls2": {},
			},
			Ecs: map[string]*Ecs{
				"testEcs1": {},
				"testEcs2": {},
			},
			Alias: map[string]*Alias{
				"testAlias1": {},
				"testAlias2": {},
			},
		}, args: args{names: []string{"testSls1", "testAlias1"}}, want: map[string]*interface{}{"testSls1": nil, "testAlias1": nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				Ecs:        tt.fields.Ecs,
				Serverless: tt.fields.Serverless,
				Alias:      tt.fields.Alias,
			}
			got := p.GetApps(tt.args.names...)
			gotKeys := maps.Keys(got)
			wantKeys := maps.Keys(tt.want)
			slices.Sort(gotKeys)
			slices.Sort(wantKeys)
			if !slices.Equal(gotKeys, wantKeys) {
				t.Errorf("GetStates() = %v, want %v", maps.Keys(got), maps.Keys(tt.want))
			}
		})
	}
}
