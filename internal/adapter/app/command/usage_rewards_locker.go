package command

import (
	"context"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/ramadani/balapan/internal/app/command"
	"github.com/ramadani/balapan/internal/app/command/model"
)

type usageRewardsLockerCommand struct {
	next   command.UsageRewardsCommander
	zkConn *zk.Conn
}

func (c *usageRewardsLockerCommand) Do(ctx context.Context, data *model.UsageRewards) error {
	locker := zk.NewLock(c.zkConn, fmt.Sprintf("/%s", data.ID), zk.WorldACL(zk.PermAll))

	if err := locker.Lock(); err != nil {
		return err
	}
	defer locker.Unlock()

	return c.next.Do(ctx, data)
}

func NewUsageRewardsLockerCommand(next command.UsageRewardsCommander, zkConn *zk.Conn) command.UsageRewardsCommander {
	return &usageRewardsLockerCommand{next: next, zkConn: zkConn}
}
