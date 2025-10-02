package delivery

import (
	"event_sourcing_bank_system_gateway/package/logger"
	"event_sourcing_bank_system_gateway/package/settings"
	"fmt"
	"net/http"
	"strconv"

	"event_sourcing_bank_system_gateway/package/grpc"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"event_sourcing_bank_system_gateway/application/model"
	"event_sourcing_bank_system_gateway/application/routing"
	"event_sourcing_bank_system_gateway/constant"
	"event_sourcing_bank_system_gateway/package/wrapper"

	"github.com/golang-jwt/jwt"
)

type RoutingHandler struct {
	config        *settings.Config
	routingUC     routing.RoutingUseCase
	registry      map[string]routingConfig
	internalRoute map[string]struct{}
}

type AuthClaims struct {
	jwt.StandardClaims
	UserID int64  `json:"Sub"`
	Email  string `json:"Email"`
}

type ResponseError struct {
	HttpCode  int    `json:"http_code"`
	GrpcCode  int    `json:"grpc_code"`
	Message   string `json:"message"`
	RootError string `json:"root_error"`
}

type AppError struct {
	Errors []ResponseError `json:"errors"`
}

func NewRoutingHandler(cfg *settings.Config, routingUC routing.RoutingUseCase) *RoutingHandler {
	return &RoutingHandler{
		config:        cfg,
		routingUC:     routingUC,
		registry:      buildRegistry(cfg),
		internalRoute: map[string]struct{}{},
	}
}

func (h *RoutingHandler) handle() gin.HandlerFunc {
	return wrapper.WithContext(func(ctx *wrapper.Context) {
		log := logger.DefaultLogger()

		route := ctx.Request.Method + ":" + ctx.FullPath()
		routingCfg, found := h.registry[route]
		if !found {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error has occurred. Please retry your request later"})
			return
		}

		uid := ctx.GetInt64(constant.KeyUserLoginID)
		metadata := map[string]string{
			"ip-address":      ctx.ClientIP(),
			"accept-language": ctx.GetHeader("Accept-Language"),
			"token":           ctx.GetString(constant.KeyAuthToken),
			"uid":             strconv.FormatInt(uid, 10),
			"user-email":      ctx.GetString(constant.KeyUserLoginEmail),
		}

		data, err := routingCfg.handler.Handle(ctx)
		if err != nil {
			appErr := &AppError{
				Errors: []ResponseError{{
					HttpCode:  http.StatusBadRequest,
					GrpcCode:  int(codes.InvalidArgument),
					Message:   err.Error(),
					RootError: err.Error(),
				}},
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, appErr)
			return
		}
		if len(routingCfg.remoteServiceName) == 0 {
			ctx.JSON(http.StatusOK, data)
			return
		}

		routing := &model.RoutingData{
			ServiceName:   routingCfg.remoteServiceName,
			ServiceMethod: routingCfg.remoteServiceMethod,
			Payload:       data,
			Metadata:      metadata,
		}

		res, err := h.routingUC.Forward(routing)
		if err != nil {
			appErr := &AppError{}
			log.Error("Forward request failed: err", zap.Error(err))
			errStatus, _ := status.FromError(err)
			for _, detail := range errStatus.Details() {
				switch info := detail.(type) {
				default:
					fmt.Printf("Unknown type: %T\n", info)
					appErr.Errors = append(appErr.Errors, ResponseError{
						GrpcCode:  int(errStatus.Code()),
						Message:   errStatus.String(),
						RootError: "",
					})
				}
			}

			if len(appErr.Errors) == 0 {
				appErr.Errors = append(appErr.Errors, ResponseError{
					GrpcCode:  int(errStatus.Code()),
					Message:   errStatus.String(),
					RootError: "",
				})
			}
			code := grpc.MapGRPCErrCodeToHttpStatus(errStatus.Code())
			ctx.AbortWithStatusJSON(code, appErr)
			return
		}

		ctx.JSON(http.StatusOK, res)
	})
}

func (h *RoutingHandler) Authorization() gin.HandlerFunc {
	return wrapper.WithContext(func(ctx *wrapper.Context) {
		ctx.Next()
	})
}
