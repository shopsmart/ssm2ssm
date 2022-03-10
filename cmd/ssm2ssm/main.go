package main

import "github.com/shopsmart/ssm2ssm/cmd"

var version = "development"

func main() {
	cmd.Execute(version, nil)
}
