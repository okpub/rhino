package actor

import (
	"fmt"
)

func decodeMiddleware(next ReceiverFunc) ReceiverFunc {
	return func(ctx ReceiverContext, data MessageEnvelope) {
		fmt.Println("链路")
		next(ctx, data)
	}
}

func uninit() {
	c1 := makeReceiverMiddlewareChain(func(ctx ReceiverContext, data MessageEnvelope) {
		fmt.Println("链路收尾")
	}, decodeMiddleware)

	c1(nil, nil)
}
