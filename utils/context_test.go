package utils_test

import (
	"context"
	"os"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/tomatorpg/tomatorpg/utils"
)

func TestWithLogger(t *testing.T) {
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	ctx := utils.WithLogger(context.Background(), logger)
	loggerOut := utils.GetLogger(ctx)
	if want, have := logger, loggerOut; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
