package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFMCService(t *testing.T) {
	ctx := NewContext(&ContextOptions{
		ENV: NewEnv(),
	})

	fcmSvc := NewFMC(ctx)
	assert.NotNil(t, fcmSvc)

	payload := &IFMCMessage{}

	ierr := fcmSvc.SendSimpleMessage([]string{"fdsdfsd"}, payload)
	assert.Error(t, ierr)
}
