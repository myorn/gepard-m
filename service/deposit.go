package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/myorn/gepard-m/constants"
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

	err = dao.AddBalanceFromTx(ctx, req.State, req.Amount)
	if err != nil {
		return &proto.Response{Message: fmt.Sprintf(MessageErr, err)}, err
	}

	err = dao.SaveTx(ctx, jsonedReq, req.TxId)
	if err != nil {
		return &proto.Response{Message: fmt.Sprintf(MessageErr, err)}, err
	}

	return &proto.Response{Message: MessageAcc}, nil
}

func validateDepositRequest(req *proto.Request) error {
	amount, err := strconv.Atoi(req.Amount)

	_, errTxId := strconv.Atoi(req.TxId)

	switch {
	case err != nil:
		return err
	case amount < 0:
		return errors.New("amount can't be negative")
	case req.Source != constants.SourceService && req.Source != constants.SourcePayment &&
		req.Source != constants.SourceGame:
		return errors.New("wrong source type")
	case req.State != constants.StateDeposit && req.State != constants.StateWithdraw:
		return errors.New("wrong state type")
	case errTxId != nil:
		return errors.New("transaction ID must be a numeric string")
	}

	return nil
}
