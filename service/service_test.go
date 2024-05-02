package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/myorn/gepard-m/constants"
	"github.com/myorn/gepard-m/db"
	"github.com/myorn/gepard-m/dto/proto"
	"github.com/myorn/gepard-m/service"
)

func contextWithDB() context.Context {
	session := db.InitDB()
	db.Flush(session)
	db.Migrate(session)
	return context.WithValue(context.Background(), constants.DBSession, session)
}

func TestPerformDepositAction_ValidRequest(t *testing.T) {
	s := service.New()
	req := &proto.Request{
		Amount: "100",
		Source: constants.SourceService,
		State:  constants.StateDeposit,
		TxId:   "123",
	}

	p, err := s.PeformDepositAction(contextWithDB(), req)

	assert.NoError(t, err)
	assert.Equal(t, service.MessageAcc, p.Message)
}

func TestPerformDepositAction_InvalidAmount(t *testing.T) {
	s := service.New()
	req := &proto.Request{
		Amount: "abc",
		Source: constants.SourceService,
		State:  constants.StateDeposit,
		TxId:   "123",
	}

	p, err := s.PeformDepositAction(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, p.Message, "Bad data")
}

func TestPerformDepositAction_NegativeAmount(t *testing.T) {
	s := service.New()
	req := &proto.Request{
		Amount: "-100",
		Source: constants.SourceService,
		State:  constants.StateDeposit,
		TxId:   "123",
	}

	p, err := s.PeformDepositAction(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, p.Message, "Bad data")
}

func TestPerformDepositAction_InvalidSource(t *testing.T) {
	s := service.New()
	req := &proto.Request{
		Amount: "100",
		Source: "invalid",
		State:  constants.StateDeposit,
		TxId:   "123",
	}

	p, err := s.PeformDepositAction(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, p.Message, "Bad data")
}

func TestPerformDepositAction_InvalidState(t *testing.T) {
	s := service.New()
	req := &proto.Request{
		Amount: "100",
		Source: constants.SourceService,
		State:  "invalid",
		TxId:   "123",
	}

	p, err := s.PeformDepositAction(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, p.Message, "Bad data")
}
