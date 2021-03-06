package interactor

import (
	"context"
	"fmt"
	"log"

	"github.com/mpppk/imagine/util"

	old0_0_1 "github.com/mpppk/imagine/domain/model/old/old0.0.1"

	"github.com/mpppk/imagine/domain/model"

	"github.com/blang/semver/v4"
	"github.com/mpppk/imagine/domain/repository"
)

type Migration struct {
	assetRepository repository.Asset
	metaRepository  repository.Meta
}

func NewMigration(assetRepository repository.Asset, metaRepository repository.Meta) *Migration {
	return &Migration{
		assetRepository: assetRepository,
		metaRepository:  metaRepository,
	}
}

func (m *Migration) Migrate() error {
	dbV, ok, err := m.metaRepository.GetDBVersion()
	if err != nil {
		return fmt.Errorf("failed to get db version: %w", err)
	}

	appV := semver.MustParse(util.Version)
	if !ok {
		if err := m.metaRepository.SetDBVersion(&appV); err != nil {
			return err
		}
		dbV = &appV
		log.Printf("info: versions: db:%s app:%s", "empty→"+appV.String(), appV.String())
	} else {
		log.Printf("info: versions: db:%s app:%s", dbV.String(), appV.String())
	}

	if err := m.migrateDB(dbV); err != nil {
		return err
	}
	return nil
}

func (m *Migration) migrateDB(dbV *semver.Version) error {
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
	c, err := m.assetRepository.ListRawByAsync(context.Background(), ws, f, 100)
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
			if _, _, err := m.assetRepository.BatchUpdateByID(ws, assets); err != nil {
				return err
			}
			log.Printf("debug: assets(ID:%v-%v) are migrated. (v0.0.1→v0.1.0)", assets[0].ID, assets[len(assets)-1].ID)
			assets = make([]*model.Asset, 0, batchNum)
		}
	}

	if _, _, err := m.assetRepository.BatchUpdateByID(ws, assets); err != nil {
		return err
	}
	if len(assets) > 0 {
		log.Printf("debug: assets(ID:%v-%v) are migrated. (v0.0.1→v0.1.0)", assets[0].ID, assets[len(assets)-1].ID)
	}

	v := semver.MustParse("0.1.0")
	return m.metaRepository.SetDBVersion(&v)
}
