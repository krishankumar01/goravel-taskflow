package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support/carbon"

	"taskflow/app/enums"
	"taskflow/app/facades"
	"taskflow/app/models"
)

type PruneExports struct {
}

func (r *PruneExports) Signature() string {
	return "exports:prune"
}

func (r *PruneExports) Description() string {
	return "Delete old exports and exports stuck in pending"
}

func (r *PruneExports) Extend() command.Extend {
	return command.Extend{
		Category: "exports",
		Flags: []command.Flag{
			&command.IntFlag{
				Name:  "days",
				Value: 7,
				Usage: "delete exports older than this many days",
			},
		},
	}
}

func (r *PruneExports) Handle(ctx console.Context) error {
	var old []models.Export
	err := facades.Orm().Query().
		Where("created_at < ?", carbon.Now().SubDays(ctx.OptionInt("days")).StdTime()).
		Find(&old)
	if err != nil {
		return err
	}
	for _, export := range old {
		if export.Path != "" && facades.Storage().Exists(export.Path) {
			if err := facades.Storage().Delete(export.Path); err != nil {
				return err
			}
		}
		if _, err := facades.Orm().Query().Delete(&export); err != nil {
			return err
		}
	}

	stuck, err := facades.Orm().Query().
		Where("status", enums.ExportPending).
		Where("created_at < ?", carbon.Now().SubHour().StdTime()).
		Delete(&models.Export{})
	if err != nil {
		return err
	}

	ctx.Info(fmt.Sprintf("Deleted %d old exports and %d stuck in pending.", len(old), stuck.RowsAffected))

	return nil
}
