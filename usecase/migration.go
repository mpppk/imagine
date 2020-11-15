package usecase

import (
	"context"
	"fmt"
	"log"

	old0_0_1 "github.com/mpppk/imagine/domain/model/old/old0.0.1"

	"github.com/mpppk/imagine/domain/model"

	"github.com/blang/semver/v4"
	"github.com/mpppk/imagine/domain/repository"
)

type Migration struct {
	client *repository.Client
}

func NewMigration(client *repository.Client) *Migration {
	return &Migration{
		client: client,
	}
}

func (m *Migration) Migrate(dbV *semver.Version) error {
	tmpl := "info: db migration is started from %s to %s"
	curDBV := dbV

	if v := semver.MustParse("0.1.0"); curDBV.Compare(v) == -1 {
		log.Printf(tmpl, curDBV, "0.1.0")
		if err := m.migrateFrom0d0d1To0d1d0("default-workspace"); err != nil {
			return fmt.Errorf("failed to migrate from 0.0.1: %w", err)
		}
		curDBV = &v
	}

	return nil
}

func (m *Migration) migrateFrom0d0d1To0d1d0(ws model.WSName) error {
	f := func(v []byte) bool {
		return true
	}
	c, err := m.client.Asset.ListRawByAsync(context.Background(), ws, f, 100)
	if err != nil {
		return err
	}
	batchNum := 10000
	var assets = make([]*model.Asset, 0, batchNum)
	for v := range c {
		oldAsset, err := old0_0_1.NewAssetFromJson(v)
		if err != nil {
			return err
		}
		asset, ok := oldAsset.Migrate()
		if !ok {
			continue
		}
		assets = append(assets, asset)
		if len(assets) >= batchNum {
			if err := m.client.Asset.BatchUpdate(ws, assets); err != nil {
				return err
			}
			log.Printf("debug: assets(ID:%v-%v) are migrated. (v0.0.1→v0.1.0)", assets[0].ID, assets[len(assets)-1].ID)
			assets = make([]*model.Asset, 0, batchNum)
		}
	}

	if err := m.client.Asset.BatchUpdate(ws, assets); err != nil {
		return err
	}
	if len(assets) > 0 {
		log.Printf("debug: assets(ID:%v-%v) are migrated. (v0.0.1→v0.1.0)", assets[0].ID, assets[len(assets)-1].ID)
	}

	v := semver.MustParse("0.1.0")
	return m.client.Meta.SetVersion(&v)
}
