package util

import (
	"testing"
)

func TestGetEncodedInfo(t *testing.T) {
	var tests = []struct {
		info, token, ans string
	}{
		{
			`{"username":"12345678","password":"123456","ip":"10.200.21.4","acid":"1","enc_ver":"srun_bx1"}`,
			"b20885b4009b8085d025c2762e3b5765c1161d8bb6dea1a48eb27bb4788dcaec",
			"{SRBX1}XzOZOusMk4UlvFoumP+nSm1SJla9zsTIi+rYt8562Xh5Y2YLr7UalkgHS/Ad3GegwekNn5rTGshz6Md8YifHMv/ThbuF0vQfuI0SCvh5AAht9SyH/FcYjaqHk8tuw0aXIcFIsS==",
		},
		{
			`{"username":"87654321","password":"654321","ip":"10.200.21.1","acid":"1","enc_ver":"srun_bx1"}`,
			"61304ea03a704dec2994577c527921f0f066c0a10e81419aaa9a77cca545fbb7",
			"{SRBX1}6h9Zmmgg3PcMdS1gYMin9KubVcnr3JVs7lubINAC9RER+YyehPUxfMtY3TBqrr1XxYrK/HZquaSXTdRkTtJxLFOh72xBwQsM8PpLGb5liZfINOZf1r9ElAIf+eyA3rCD48AxeS==",
		},
	}

	for _, val := range tests {
		get := GetEncodedInfo(val.info, val.token)
		if get != val.ans {
			t.Errorf("\nExcepted:\n%s\nWhile found:\n%s\n", val.ans, get)
		}
	}
}
