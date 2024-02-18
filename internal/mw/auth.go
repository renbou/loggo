package mw

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
)

const (
	pigeonAuthHeader = "X-Loggo-Pigeon"
	pigeonUnknown    = "_unknown_"
)

type pigeonNameCtxKey struct{}

// PigeonAuthStreamServerInteceptor returns a new StreamServerInterceptor which authenticates requests
// from pigeons using a token->name map provider on the specified method prefix.
// Note that pigeonProvider will be called concurrently on each request to receive new tokens once they are added.
func PigeonAuthStreamServerInteceptor(pigeonProvider func() map[string]string, prefix string,
) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if strings.HasPrefix(info.FullMethod, prefix) {
			tokens := pigeonProvider()

			headers := metadata.ValueFromIncomingContext(ss.Context(), pigeonAuthHeader)
			if len(headers) != 1 {
				return status.Error(codes.Unauthenticated, "pigeon token not specified properly")
			}

			// Getting a value from a map hashes it first, so this should be resistant to timing attacks
			name, ok := tokens[headers[0]]
			if !ok {
				return status.Error(codes.Unauthenticated, "invalid pigeon token")
			}

			ss = pigeonNameToServerStream(ss, name)
		}

		return handler(srv, ss)
	}
}

func pigeonNameToServerStream(ss grpc.ServerStream, name string) grpc.ServerStream {
	ws := grpc_middleware.WrapServerStream(ss)
	ws.WrappedContext = context.WithValue(ss.Context(), pigeonNameCtxKey{}, name)

	return ws
}

// PigeonNameFromCtx gets the pigeon name set in PigeonAuth interceptors.
func PigeonNameFromCtx(ctx context.Context) string {
	value := ctx.Value(pigeonNameCtxKey{})
	name, ok := value.(string)

	if !ok {
		return pigeonUnknown
	}
	return name
}
