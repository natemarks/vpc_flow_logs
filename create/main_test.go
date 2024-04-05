package create

import (
	"testing"

	"github.com/natemarks/vpc_flow_logs/config"
	"github.com/rs/zerolog"
)

func TestFlowLog(t *testing.T) {
	t.Skip("Skipping test - don't need to hit AWS all the time")
	myConfig, err := config.GetCreateConfig("vpc-0d354f2e35a217375", true)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}
	logger := config.GetLogger(false)
	type args struct {
		cConfig config.CreateConfig
		log     *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test1",
			args: args{
				cConfig: myConfig,
				log:     &logger,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FlowLog(tt.args.cConfig, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlowLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
