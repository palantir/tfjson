/*
Copyright (c) 2016 Palantir Technologies

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/terraform/terraform"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: tfjson terraform.tfplan")
		os.Exit(1)
	}

	j, err := tfjson(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(j)
}

type output map[string]interface{}

func tfjson(planfile string) (string, error) {
	f, err := os.Open(planfile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return "", err
	}

	diff := output{}
	for _, v := range plan.Diff.Modules {
		convertModuleDiff(diff, v)
	}

	j, err := json.MarshalIndent(diff, "", "    ")
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func insert(out output, path []string, key string, value interface{}) {
	if len(path) > 0 && path[0] == "root" {
		path = path[1:]
	}
	for _, elem := range path {
		switch nested := out[elem].(type) {
		case output:
			out = nested
		default:
			new := output{}
			out[elem] = new
			out = new
		}
	}
	out[key] = value
}

func convertModuleDiff(out output, diff *terraform.ModuleDiff) {
	insert(out, diff.Path, "destroy", diff.Destroy)
	for k, v := range diff.Resources {
		convertInstanceDiff(out, append(diff.Path, k), v)
	}
}

func convertInstanceDiff(out output, path []string, diff *terraform.InstanceDiff) {
	insert(out, path, "destroy", diff.Destroy)
	insert(out, path, "destroy_tainted", diff.DestroyTainted)
	for k, v := range diff.Attributes {
		insert(out, path, k, v.New)
	}
}
