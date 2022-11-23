package mq

import (
	"context"
)

type ConsumeCallBack func(ctx context.Context, msgId string, msg []byte, extra interface{}) error

type Contracts interface {
	SendNormalMsg(ctx context.Context, data []byte) (err error, msgId string)
	ReceiveNormalMsg(ctx context.Context, fn ConsumeCallBack)

	SendDelayMsg(ctx context.Context, data []byte, delay int64) (err error, msgId string)
	ReceiveDelayMsg(ctx context.Context, fn ConsumeCallBack)
}

type Service struct {
	Driver Contracts
}

func NewService(driver Contracts) *Service {
	return &Service{
		Driver: driver,
	}
}

func (s *Service) SendNormalMsg(ctx context.Context, data []byte) (err error, msgId string) {
	return s.Driver.SendNormalMsg(ctx, data)
}

func (s *Service) ReceiveNormalMsg(ctx context.Context, fn ConsumeCallBack) {
	s.Driver.ReceiveNormalMsg(ctx, fn)
}

func (s *Service) SendDelayMsg(ctx context.Context, data []byte, delay int64) (err error, msgId string) {
	if delay <= 0 {
		panic("SendDelayMsg must provide delay")
	}
	return s.Driver.SendDelayMsg(ctx, data, delay)
}

func (s *Service) ReceiveDelayMsg(ctx context.Context, fn ConsumeCallBack) {
	s.Driver.ReceiveDelayMsg(ctx, fn)
}
