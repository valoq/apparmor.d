#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=systemd-sysusers
@test "systemd-sysusers: Print the contents of all configuration files (before each file, its name is printed as a comment)" {
    systemd-sysusers --cat-config
    aa_check
}

# bats test_tags=systemd-sysusers
@test "systemd-sysusers: Process configuration files and print what would be done without actually doing anything" {
    systemd-sysusers --dry-run
    aa_check
}

# bats test_tags=systemd-sysusers
@test "systemd-sysusers: Create users and groups from all configuration file" {
    sudo systemd-sysusers
    aa_check
}