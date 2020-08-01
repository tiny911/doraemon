package log

import "testing"

var (
	path = "./"
	prj  = "testPrj"
)

func TestLogHelper_Execute(t *testing.T) {
	helper := NewLogHelper(prj, path)
	helper.Execute()
	helper.Cancel()
	t.Logf("execute success.")
}
