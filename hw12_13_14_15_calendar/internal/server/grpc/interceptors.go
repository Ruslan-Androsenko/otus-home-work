package grpc

import (
	"context"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/server"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	protocolField         = "protocol"
	grpcServiceField      = "grpc.service"
	grpcMethodField       = "grpc.method"
	peerAddressField      = "peer.address"
	grpcSendDurationField = "grpc.send.duration"
)

func InterceptorLogger(l server.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		var protocolValue, grpcServiceValue, grpcMethodValue, peerAddressValue, grpcSendDurationValue string

		// Производим маппинг полей и их значений
		for i := 0; i < len(fields); i += 2 {
			fieldName, okName := fields[i].(string)
			fieldValue, okValue := fields[i+1].(string)

			if okName && okValue {
				switch fieldName {
				case protocolField:
					protocolValue = fieldValue

				case grpcServiceField:
					grpcServiceValue = fieldValue

				case grpcMethodField:
					grpcMethodValue = fieldValue

				case peerAddressField:
					peerAddressValue = fieldValue

				case grpcSendDurationField:
					grpcSendDurationValue = fieldValue
				}
			}
		}

		l.Infof("%s %s/%s %s %s",
			peerAddressValue, grpcServiceValue, grpcMethodValue,
			protocolValue, grpcSendDurationValue,
		)
	})
}
