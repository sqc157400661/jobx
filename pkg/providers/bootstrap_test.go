package providers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClone(t *testing.T) {
	ob1 := &DelayTasker{}
	Set(ob1)
	ob2 := Get(ob1.Name())
	assert.NotSame(t, ob1, ob2)
	assert.IsType(t, ob1, ob2)
	assert.Equal(t, true, assert.ObjectsAreEqual(ob1, ob2))
	//fmt.Printf("origin:%p  clone:%p \n", ob1, ob2)
	//fmt.Printf("origin:%T  clone:%T \n", ob1, ob2)
	//fmt.Printf("origin:%+v  clone:%+v \n", ob1, ob2)
}
