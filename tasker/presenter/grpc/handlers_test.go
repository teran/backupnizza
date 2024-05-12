//go:build grpc

package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	"github.com/teran/backupnizza/tasker/presenter/grpc/proto"
	"github.com/teran/backupnizza/tasker/service"
	"github.com/teran/go-grpctest"
)

func (s *handlersTestSuite) TestRunTask() {
	s.serviceMock.On("GetByName", "test task").Return(&testRunner{}, nil).Once()
	_, err := s.client.RunTask(s.ctx, &proto.RunTaskRequest{
		Name: "test task",
	})
	s.Require().NoError(err)
}

func (s *handlersTestSuite) TestRunTaskWithError() {
	s.serviceMock.On("GetByName", "test task").Return(&testRunner{
		err: errors.New("test error"),
	}, nil).Once()
	_, err := s.client.RunTask(s.ctx, &proto.RunTaskRequest{
		Name: "test task",
	})
	s.Require().Error(err)
	s.Require().Equal("rpc error: code = Internal desc = test error", err.Error())
}

// ========================================================================
// Test suite setup
// ========================================================================
type handlersTestSuite struct {
	suite.Suite

	serviceMock *service.Mock

	srv    grpctest.Server
	ctx    context.Context
	cancel context.CancelFunc

	client   proto.TaskerServiceClient
	handlers Handlers
}

func (s *handlersTestSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 30*time.Second)

	s.serviceMock = service.NewMock()

	s.handlers = New(s.serviceMock)

	s.srv = grpctest.New()
	s.handlers.Register(s.srv.Server())

	err := s.srv.Run()
	s.Require().NoError(err)

	dial, err := s.srv.DialContext(s.ctx)
	s.Require().NoError(err)

	s.client = proto.NewTaskerServiceClient(dial)
}

func (s *handlersTestSuite) TearDownTest() {
	s.srv.Close()
	s.cancel()
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, &handlersTestSuite{})
}

type testRunner struct {
	err error
}

func (tr *testRunner) Run(ctx context.Context) error {
	return tr.err
}
