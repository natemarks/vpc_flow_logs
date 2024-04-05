package delete

import (
	"testing"

	"github.com/natemarks/vpc_flow_logs/config"
	"github.com/rs/zerolog"
)

func TestFlowLog(t *testing.T) {
	t.Skip("skipping test") // skipping DELETE test
	log := config.GetLogger(true)
	deleteConfig, err := config.GetDeleteConfig("flow-log-id", true)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		cfg config.DeleteConfig
		log *zerolog.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test FlowLog",
			args: args{
				cfg: deleteConfig,
				log: &log,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FlowLog(tt.args.cfg, tt.args.log)
		})
	}
}
