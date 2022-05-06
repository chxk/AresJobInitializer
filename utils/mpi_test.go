package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteMPIHostFile(t *testing.T) {
	os.Setenv("LOCAL_IP", "localhost")
	ips := []string{"1.2.3.4", "5.6.7.8"}
	f := "./hostfile,./hostfile1"
	err := WriteMPIHostFile(OMPI, f, ips, 2)
	assert.Equal(t, err, nil)

	err = WriteMPIHostFile(MPICH, f, ips, 2)
	assert.Equal(t, err, nil)

	err = WriteMPIHostFile("", f, ips, 2)
	assert.Equal(t, err, nil)

	os.Unsetenv("LOCAL_IP")
	err = WriteMPIHostFile("", f, nil, 2)
	assert.NotNil(t, err)
}
