// aa-log - Review AppArmor generated messages
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"path/filepath"
	"testing"
)

var (
	testdata = "../../tests/testdata/logs"
)

func Test_app(t *testing.T) {
	tests := []struct {
		name    string
		logger  string
		path    string
		profile string
		rules   bool
		wantErr bool
	}{
		{
			name:    "Test audit.log",
			logger:  "auditd",
			path:    filepath.Join(testdata, "audit.log"),
			profile: "",
			rules:   false,
			wantErr: false,
		},
		{
			name:    "Test audit.log to rules",
			logger:  "auditd",
			path:    filepath.Join(testdata, "audit.log"),
			profile: "",
			rules:   true,
			wantErr: false,
		},
		{
			name:    "Test Dbus Session",
			logger:  "systemd",
			path:    filepath.Join(testdata, "systemd.log"),
			profile: "",
			rules:   false,
			wantErr: false,
		},
		{
			name:    "No logfile",
			logger:  "auditd",
			path:    filepath.Join(testdata, "log"),
			profile: "",
			rules:   false,
			wantErr: true,
		},
		{
			name:    "Logger not supported",
			logger:  "raw",
			path:    filepath.Join(testdata, "audit.log"),
			profile: "",
			rules:   false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules = tt.rules
			if err := aaLog(tt.logger, tt.path, tt.profile); (err != nil) != tt.wantErr {
				t.Errorf("aaLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
