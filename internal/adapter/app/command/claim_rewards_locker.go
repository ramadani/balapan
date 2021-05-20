package command

import (
	"context"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/ramadani/balapan/internal/app/command"
	"github.com/ramadani/balapan/internal/app/command/model"
)

type claimRewardsLockerCommand struct {
	next   command.ClaimRewardsCommander
	zkConn *zk.Conn
}

func (c *claimRewardsLockerCommand) Do(ctx context.Context, data *model.ClaimRewards) error {
	locker := zk.NewLock(c.zkConn, fmt.Sprintf("/%s", data.ID), zk.WorldACL(zk.PermAll))

	if err := locker.Lock(); err != nil {
		return err
	}
	defer locker.Unlock()

	return c.next.Do(ctx, data)
}

func NewClaimRewardsLockerCommand(next command.ClaimRewardsCommander, zkConn *zk.Conn) command.ClaimRewardsCommander {
	return &claimRewardsLockerCommand{next: next, zkConn: zkConn}
}
