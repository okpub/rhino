package ini

import (
	"testing"
)

//test
func TestINI(t *testing.T) {
	Unmarshal("./server.conf")
}
