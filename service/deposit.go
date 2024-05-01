package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/myorn/gepard-m/dao"
	"github.com/myorn/gepard-m/dto/proto"
)

type depositServer struct {
	proto.UnimplementedDepositServer
}

const (
	MessageAcc    = "Accepted"
	MessageErr    = "Error: %v"
	MessageBadReq = "Bad data: %v"

	// accepted sources
	sourceGame    = "game"
	sourcePayment = "payment"
	sourceService = "service"

	// accepted states
	stateDeposit  = "deposit"
	stateWithdraw = "withdraw"
)

func New() *depositServer {
	return &depositServer{}
}

func (*depositServer) PeformDepositAction(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	// validate
	err := validateDepositRequest(req)
	if err != nil {
		return &proto.Response{Message: fmt.Sprintf(MessageBadReq, err)}, err
	}

	jsonedReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = dao.SaveTxAndUpdateBalance(ctx, jsonedReq, req.Amount, req.State, req.TxId)
	if err != nil {
		return &proto.Response{Message: fmt.Sprintf(MessageErr, err)}, err
	}

	return &proto.Response{Message: MessageAcc}, nil
}

func validateDepositRequest(req *proto.Request) error {
	amount, err := strconv.Atoi(req.Amount)

	switch {
	case err != nil:
		return err
	case amount < 0:
		return errors.New("amount can't be negative")
	case req.Source != sourceService && req.Source != sourcePayment &&
		req.Source != sourceGame:
		return errors.New("wrong source type")
	case req.State != stateDeposit && req.State != stateWithdraw:
		return errors.New("wrong state type")
	}

	return nil
}
