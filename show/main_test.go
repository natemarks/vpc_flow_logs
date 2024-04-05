package show

import (
	"testing"

	"github.com/natemarks/vpc_flow_logs/types"

	"github.com/rs/zerolog"
)

func TestGetFlowLogDescriptions(t *testing.T) {
	t.Skip("Skipping test - don't need to hit AWS all the time")
	type args struct {
		logger *zerolog.Logger
	}
	tests := []struct {
		name         string
		args         args
		wantFlowLogs []types.FlowLog
		wantErr      bool
	}{
		{
			name: "Test GetFlowLogDescriptions",
			args: args{
				logger: &zerolog.Logger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetFlowLogDescriptions(tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFlowLogDescriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
