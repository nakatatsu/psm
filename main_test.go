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
			name: "sync basic",
			args: []string{"psm", "sync", "file.yaml"},
			want: Config{Subcommand: "sync", File: "file.yaml"},
		},
		{
			name: "sync with dry-run",
			args: []string{"psm", "sync", "--dry-run", "file.yaml"},
			want: Config{Subcommand: "sync", DryRun: true, File: "file.yaml"},
		},
		{
			name: "sync with skip-approve",
			args: []string{"psm", "sync", "--skip-approve", "file.yaml"},
			want: Config{Subcommand: "sync", SkipApprove: true, File: "file.yaml"},
		},
		{
			name: "sync with debug",
			args: []string{"psm", "sync", "--debug", "file.yaml"},
			want: Config{Subcommand: "sync", Debug: true, File: "file.yaml"},
		},
		{
			name: "sync with delete file",
			args: []string{"psm", "sync", "--delete", "needless.yml", "file.yaml"},
			want: Config{Subcommand: "sync", DeleteFile: "needless.yml", File: "file.yaml"},
		},
		{
			name: "sync all flags combined",
			args: []string{"psm", "sync", "--dry-run", "--skip-approve", "--debug", "--delete", "del.yml", "file.yaml"},
			want: Config{Subcommand: "sync", DryRun: true, SkipApprove: true, Debug: true, DeleteFile: "del.yml", File: "file.yaml"},
		},
		{
			name: "export basic",
			args: []string{"psm", "export", "out.yaml"},
			want: Config{Subcommand: "export", File: "out.yaml"},
		},
		{
			name: "export with debug",
			args: []string{"psm", "export", "--debug", "out.yaml"},
			want: Config{Subcommand: "export", Debug: true, File: "out.yaml"},
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
			name:    "sync no file arg",
			args:    []string{"psm", "sync"},
			wantErr: true,
		},
		{
			name:    "export no file arg",
			args:    []string{"psm", "export"},
			wantErr: true,
		},
		{
			name:    "prune flag removed",
			args:    []string{"psm", "sync", "--prune", "file.yaml"},
			wantErr: true,
		},
		{
			name:    "store flag removed with ssm",
			args:    []string{"psm", "sync", "--store", "ssm", "file.yaml"},
			wantErr: true,
		},
		{
			name:    "store flag removed with sm",
			args:    []string{"psm", "export", "--store", "sm", "out.yaml"},
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
