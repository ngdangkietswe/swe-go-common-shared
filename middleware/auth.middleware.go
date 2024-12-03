package middleware

import (
	"context"
	"fmt"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/ngdangkietswe/swe-go-common-shared/grpc/constant"
	grpcutil "github.com/ngdangkietswe/swe-go-common-shared/grpc/util"
	"github.com/ngdangkietswe/swe-go-common-shared/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	JwtBearerPrefix     = "Bearer "
	AuthorizationHeader = "authorization"
)

// AuthMiddleware is a middleware function that checks the token in the request header
func AuthMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata not found")
	}

	tokens := md.Get(AuthorizationHeader)
	if len(tokens) > 0 {
		value := tokens[0]
		if value == "" || !strings.HasPrefix(value, JwtBearerPrefix) {
			return nil, fmt.Errorf("missing or invalid token")
		} else {
			token := strings.TrimPrefix(value, JwtBearerPrefix)
			jwtClaims, err := util.ParseToken(token, config.GetString("JWT_SECRET", ""))
			if err != nil {
				return nil, fmt.Errorf("invalid token")
			}

			principal, err := grpcutil.AsGrpcPrincipal(jwtClaims)
			if err != nil {
				return nil, fmt.Errorf("invalid token")
			}

			ctx = context.WithValue(ctx, constant.CtxPrincipalKey, principal)
		}
	}

	return handler(ctx, req)
}
