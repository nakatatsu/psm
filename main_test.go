package main

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Config
		wantErr bool
	}{
		{
			name: "sync with ssm",
			args: []string{"psm", "sync", "--store", "ssm", "file.yaml"},
			want: Config{Subcommand: "sync", Store: "ssm", File: "file.yaml"},
		},
		{
			name:    "sync with sm is invalid",
			args:    []string{"psm", "sync", "--store", "sm", "--profile", "prod", "file.yaml"},
			wantErr: true,
		},
		{
			name: "sync with dry-run",
			args: []string{"psm", "sync", "--store", "ssm", "--dry-run", "file.yaml"},
			want: Config{Subcommand: "sync", Store: "ssm", DryRun: true, File: "file.yaml"},
		},
		{
			name: "sync with skip-approve",
			args: []string{"psm", "sync", "--store", "ssm", "--skip-approve", "file.yaml"},
			want: Config{Subcommand: "sync", Store: "ssm", SkipApprove: true, File: "file.yaml"},
		},
		{
			name: "sync with debug",
			args: []string{"psm", "sync", "--store", "ssm", "--debug", "file.yaml"},
			want: Config{Subcommand: "sync", Store: "ssm", Debug: true, File: "file.yaml"},
		},
		{
			name: "sync with delete file",
			args: []string{"psm", "sync", "--store", "ssm", "--delete", "needless.yml", "file.yaml"},
			want: Config{Subcommand: "sync", Store: "ssm", DeleteFile: "needless.yml", File: "file.yaml"},
		},
		{
			name: "sync all flags combined",
			args: []string{"psm", "sync", "--store", "ssm", "--dry-run", "--skip-approve", "--debug", "--delete", "del.yml", "file.yaml"},
			want: Config{Subcommand: "sync", Store: "ssm", DryRun: true, SkipApprove: true, Debug: true, DeleteFile: "del.yml", File: "file.yaml"},
		},
		{
			name: "export with ssm",
			args: []string{"psm", "export", "--store", "ssm", "out.yaml"},
			want: Config{Subcommand: "export", Store: "ssm", File: "out.yaml"},
		},
		{
			name:    "export with sm is invalid",
			args:    []string{"psm", "export", "--store", "sm", "--profile", "staging", "out.yaml"},
			wantErr: true,
		},
		{
			name: "export with debug",
			args: []string{"psm", "export", "--store", "ssm", "--debug", "out.yaml"},
			want: Config{Subcommand: "export", Store: "ssm", Debug: true, File: "out.yaml"},
		},
		{
			name:    "no subcommand",
			args:    []string{"psm"},
			wantErr: true,
		},
		{
			name:    "invalid subcommand",
			args:    []string{"psm", "invalid"},
			wantErr: true,
		},
		{
			name:    "sync no store flag",
			args:    []string{"psm", "sync", "file.yaml"},
			wantErr: true,
		},
		{
			name:    "sync invalid store value",
			args:    []string{"psm", "sync", "--store", "dynamodb", "file.yaml"},
			wantErr: true,
		},
		{
			name:    "sync no file arg",
			args:    []string{"psm", "sync", "--store", "ssm"},
			wantErr: true,
		},
		{
			name:    "export no file arg",
			args:    []string{"psm", "export", "--store", "ssm"},
			wantErr: true,
		},
		{
			name:    "export no store flag",
			args:    []string{"psm", "export", "out.yaml"},
			wantErr: true,
		},
		{
			name:    "prune flag removed",
			args:    []string{"psm", "sync", "--store", "ssm", "--prune", "file.yaml"},
			wantErr: true,
		},
		{
			name: "version flag",
			args: []string{"psm", "--version"},
			want: Config{ShowVersion: true},
		},
		{
			name: "version flag with extra args",
			args: []string{"psm", "--version", "sync"},
			want: Config{ShowVersion: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArgs() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
