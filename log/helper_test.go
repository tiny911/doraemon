package log

import "testing"

var (
	path = "/data/project/"
	prj  = "testPrj"
)

func TestLogHelper_Execute(t *testing.T) {
	helper := NewLogHelper(prj, path)
	helper.Execute()
	helper.Cancel()
	t.Logf("execute success.")
}
