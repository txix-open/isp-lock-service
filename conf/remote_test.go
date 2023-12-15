package conf_test

import (
	"isp-lock-service/conf"
	"testing"

	"github.com/integration-system/isp-kit/test/rct"
)

func TestDefaultRemoteConfig(t *testing.T) {
	t.Parallel()
	rct.Test(t, "default_remote_config.json", conf.Remote{})
}
