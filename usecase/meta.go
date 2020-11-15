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
	assetRepository repository.Asset
	metaRepository  repository.Meta
}

func NewMigration(metaRepository repository.Meta) *Migration {
	return &Migration{
		metaRepository: metaRepository,
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
	// TODO workspace名を取得
	f := func(v []byte) bool {
		return true
	}
	c, err := m.assetRepository.ListRawByAsync(context.Background(), ws, f, 100)
	if err != nil {
		return err
	}

	for v := range c {
		oldAsset, err := old0_0_1.NewAssetFromJson(v)
		if err != nil {
			return err
		}
		asset := oldAsset.Migrate()
		if err := m.assetRepository.Update(ws, asset); err != nil {
			return err
		}
	}
	v := semver.MustParse("0.1.0")
	return m.metaRepository.SetVersion(&v)
}
