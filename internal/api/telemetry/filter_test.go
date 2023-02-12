package telemetry

import (
	"testing"

	"github.com/renbou/loggo/internal/storage"
	pb "github.com/renbou/loggo/pkg/api/telemetry"
	"github.com/stretchr/testify/assert"
)

func flatMappingFromMap(m map[string]string) storage.FlatMapping {
	return func(key string) (value string, ok bool) {
		v, ok := m[key]
		return v, ok
	}
}

func applyFilter(filter storage.Filter, m storage.Message, fm storage.FlatMapping) bool {
	if filter == nil {
		return true
	}

	return filter(m, fm)
}

func Test_CompileFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		raw            *pb.LogFilter
		message        string
		flatmap        map[string]string
		expectedResult bool
	}{
		{
			name:           "nil filter",
			raw:            &pb.LogFilter{Filter: nil},
			message:        `{"some":"message"}`,
			flatmap:        map[string]string{"some": "message"},
			expectedResult: true,
		},
		{
			name:           "text filter found",
			raw:            &pb.LogFilter{Filter: &pb.LogFilter_Text_{Text: &pb.LogFilter_Text{Value: "message"}}},
			message:        `{"some":"message"}`,
			flatmap:        map[string]string{"some": "message"},
			expectedResult: true,
		},
		{
			name:           "text filter not found",
			raw:            &pb.LogFilter{Filter: &pb.LogFilter_Text_{Text: &pb.LogFilter_Text{Value: "invalid"}}},
			message:        `{}`,
			flatmap:        map[string]string{},
			expectedResult: false,
		},
		{
			name:           "scoped filter found",
			raw:            &pb.LogFilter{Filter: &pb.LogFilter_Scoped_{Scoped: &pb.LogFilter_Scoped{Field: "scoped.field", Value: "1337"}}},
			message:        `{"scoped":{"field": 1337}}`,
			flatmap:        map[string]string{"scoped.field": "1337"},
			expectedResult: true,
		},
		{
			name:           "scoped field not found",
			raw:            &pb.LogFilter{Filter: &pb.LogFilter_Scoped_{Scoped: &pb.LogFilter_Scoped{Field: "unknown", Value: ""}}},
			message:        `{"different":"field"}`,
			flatmap:        map[string]string{"different": "field"},
			expectedResult: false,
		},
		{
			name: "composition",
			raw: &pb.LogFilter{Filter: &pb.LogFilter_And_{And: &pb.LogFilter_And{
				A: &pb.LogFilter{Filter: &pb.LogFilter_Or_{Or: &pb.LogFilter_Or{
					A: &pb.LogFilter{Filter: &pb.LogFilter_Text_{Text: &pb.LogFilter_Text{Value: "text value"}}},
					B: &pb.LogFilter{Filter: &pb.LogFilter_Not_{Not: &pb.LogFilter_Not{
						A: &pb.LogFilter{Filter: &pb.LogFilter_Scoped_{Scoped: &pb.LogFilter_Scoped{Field: "scoped", Value: "field"}}},
					}}},
				}}},
				B: &pb.LogFilter{Filter: &pb.LogFilter_Scoped_{Scoped: &pb.LogFilter_Scoped{Field: "other.0", Value: "field"}}},
			}}},
			message:        `{"info":"wrong value", "scoped":false, "other":["field"]}`,
			flatmap:        map[string]string{"info": "text value", "scoped": "false", "other.0": "field"},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			filter := compileFilter(tt.raw)
			gotResult := applyFilter(filter, storage.Message(tt.message), flatMappingFromMap(tt.flatmap))

			assert.Equal(t, tt.expectedResult, gotResult)
		})
	}
}
