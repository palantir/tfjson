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

	err := tfjson(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type output map[string]interface{}

func tfjson(planfile string) error {
	f, err := os.Open(planfile)
	if err != nil {
		return err
	}
	defer f.Close()

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return err
	}

	diff := output{}
	for _, v := range plan.Diff.Modules {
		convertModuleDiff(diff, v)
	}

	j, err := json.MarshalIndent(diff, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(j))
	return nil
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
