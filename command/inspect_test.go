package command

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_commands(t *testing.T) {

	tc := []struct {
		command  []string
		env      []string
		expected string
	}{
		{[]string{"inspect", "-var=fruit=banana", filepath.Join(testFixture("var-arg"), "fruit_builder.pkr.hcl")}, nil, `Packer Inspect: HCL2 mode

> input-variables:

var.fruit: "banana" [debug: {Type:cty.String,CmdValue:banana,VarfileValue:null,EnvValue:null,DefaultValue:null}]

> local-variables:

local.fruit: "banana"

> builds:

  > <unnamed build 0>:

    sources:

      null.builder

    provisioners:

      shell-local

    post-processors:

      <no post-processor>

`},
		{[]string{"inspect", "-var=fruit=peach", filepath.Join(testFixture("hcl"), "inspect", "fruit_string.pkr.hcl")}, nil, `Packer Inspect: HCL2 mode

> input-variables:

var.fruit: "peach" [debug: {Type:cty.String,CmdValue:peach,VarfileValue:null,EnvValue:null,DefaultValue:banana}]

> local-variables:


> builds:

`},
		{[]string{"inspect", "-var=fruit=peach", filepath.Join(testFixture("hcl"), "inspect")}, nil, `Packer Inspect: HCL2 mode

> input-variables:

var.fruit: "peach" [debug: {Type:cty.String,CmdValue:peach,VarfileValue:null,EnvValue:null,DefaultValue:banana}]

> local-variables:


> builds:

  > aws_example_builder:

  > Description: The builder of clouds !!

Use it at will.


    sources:

      amazon-ebs.example-1

      amazon-ebs.example-2

    provisioners:

      shell

    post-processors:

      0:
        manifest

      1:
        shell-local

      2:
        manifest
        shell-local

`},
		{[]string{"inspect", "-var=fruit=peach", filepath.Join(testFixture("var-arg"), "fruit_builder.json")}, nil, `Packer Inspect: JSON mode
Required variables:

  fruit

Optional variables and their defaults:


Builders:

  null

Provisioners:

  <No provisioners>

Note: If your build names contain user variables or template
functions such as 'timestamp', these are processed at build time,
and therefore only show in their raw form here.
`},
	}

	for _, tc := range tc {
		t.Run(fmt.Sprintf("packer %s", tc.command), func(t *testing.T) {
			p := helperCommand(t, tc.command...)
			p.Env = append(p.Env, tc.env...)
			bs, err := p.Output()
			if err != nil {
				t.Fatalf("%v: %s", err, bs)
			}
			actual := string(bs)
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Fatalf("unexpected output %s", diff)
			}
		})
	}
}
