package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/ngdangkietswe/swe-go-common-shared/constants"
	"github.com/ngdangkietswe/swe-go-common-shared/domain"
	"github.com/ngdangkietswe/swe-go-common-shared/grpc/constant"
	grpcutil "github.com/ngdangkietswe/swe-go-common-shared/grpc/util"
	"github.com/ngdangkietswe/swe-go-common-shared/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

// AuthMiddleware is a middleware function that checks the token in the request header
func AuthMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata not found")
	}

	tokens := md.Get(strings.ToLower(constants.AuthorizationHeader))
	userPermissions := md.Get(strings.ToLower(constants.GrpcMetadataUserPermission))

	if len(tokens) > 0 {
		value := tokens[0]
		if value == "" || !strings.HasPrefix(value, constants.TokenPrefix) {
			return nil, fmt.Errorf("missing or invalid token")
		} else {
			token := strings.TrimSpace(strings.TrimPrefix(value, constants.TokenPrefix))
			jwtClaims, err := util.ParseToken(token, config.GetString("JWT_SECRET", ""))
			if err != nil {
				return nil, fmt.Errorf("invalid token")
			}

			principal, err := grpcutil.AsGrpcPrincipal(jwtClaims)
			if err != nil {
				return nil, fmt.Errorf("invalid token")
			}

			if userPermissions != nil && len(userPermissions) > 0 {
				var userPermission *domain.UserPermission
				err = json.Unmarshal([]byte(userPermissions[0]), &userPermission)
				if err != nil {
					return nil, fmt.Errorf("invalid user permissions")
				}
				principal.UserPermission = userPermission
			}

			ctx = context.WithValue(ctx, constant.CtxPrincipalKey, principal)
		}
	}

	return handler(ctx, req)
}
