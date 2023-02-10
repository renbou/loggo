package mw

import (
	"context"
	"testing"

	"github.com/renbou/loggo/internal/mw/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Test_PigeonAuthStreamServerInteceptor(t *testing.T) {
	t.Parallel()

	const prefix = "/test/NeedAuth"
	const authMethod = "/test/NeedAuth/method"
	pigeonProvider := func() map[string]string {
		return map[string]string{
			"111-222-333": "pigeon",
		}
	}

	tests := []struct {
		name           string
		method         string
		md             metadata.MD
		expectedPigeon string
		wantErr        bool
	}{
		{
			name:           "auth endpoint with valid token",
			method:         authMethod,
			md:             metadata.MD{pigeonAuthHeader: []string{"111-222-333"}},
			expectedPigeon: "pigeon",
		},
		{
			name:    "auth endpoint with no token",
			method:  authMethod,
			md:      metadata.MD{},
			wantErr: true,
		},
		{
			name:    "auth endpoint with invalid token",
			method:  authMethod,
			md:      metadata.MD{pigeonAuthHeader: []string{"bad-token"}},
			wantErr: true,
		},
		{
			name:           "no auth",
			method:         "/test/NoAuth/method",
			md:             metadata.MD{pigeonAuthHeader: []string{"111-222-333"}},
			expectedPigeon: pigeonUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// No need to mock other methods since the interceptor simply works with the context
			mockStream := mocks.NewServerStreamMock(t)
			mockStream.ContextMock.Expect().Return(metadata.NewIncomingContext(context.Background(), tt.md))

			var gotPigeon string
			handler := func(srv interface{}, stream grpc.ServerStream) error {
				gotPigeon = PigeonNameFromCtx(stream.Context())
				return nil
			}

			// Act. Create the interceptor and "intercept" a request
			interceptor := PigeonAuthStreamServerInteceptor(pigeonProvider, prefix)
			err := interceptor(nil, mockStream, &grpc.StreamServerInfo{FullMethod: tt.method}, handler)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPigeon, gotPigeon)
			}
		})
	}
}
