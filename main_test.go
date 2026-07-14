package main

import "testing"

func TestParseStartupMode(t *testing.T) {
	tests := []struct {
		name string
		args []string
		cli  bool
		dir  string
		fail bool
	}{
		{"menu-default", nil, false, "", false},
		{"menu-directory", []string{"posts"}, false, "posts", false},
		{"cli-default-directory", []string{"--process-new"}, true, "", false},
		{"cli-directory", []string{"--process-new", "posts"}, true, "posts", false},
		{"too-many-arguments", []string{"--process-new", "posts", "more"}, false, "", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cli, dir, err := parseStartupMode(test.args)
			if (err != nil) != test.fail || (!test.fail && (cli != test.cli || dir != test.dir)) {
				t.Fatalf("cli=%v dir=%q err=%v", cli, dir, err)
			}
		})
	}
}
