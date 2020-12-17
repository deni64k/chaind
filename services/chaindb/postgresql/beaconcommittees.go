// Copyright © 2020 Weald Technology Trading.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgresql

import (
	"context"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/wealdtech/chaind/services/chaindb"
)

// SetBeaconCommittee sets a beacon committee.
func (s *Service) SetBeaconCommittee(ctx context.Context, beaconCommittee *chaindb.BeaconCommittee) error {
	tx := s.tx(ctx)
	if tx == nil {
		return ErrNoTransaction
	}

	_, err := tx.Exec(ctx, `
      INSERT INTO t_beacon_committees(f_slot
                                     ,f_index
                                     ,f_committee)
      VALUES($1,$2,$3)
      ON CONFLICT (f_slot,f_index) DO
      UPDATE
      SET f_committee = excluded.f_committee
		 `,
		beaconCommittee.Slot,
		beaconCommittee.Index,
		beaconCommittee.Committee,
	)

	return err
}

// BeaconCommitteeBySlotAndIndex fetches the beacon committee with the given slot and index.
func (s *Service) BeaconCommitteeBySlotAndIndex(ctx context.Context, slot spec.Slot, index spec.CommitteeIndex) (*chaindb.BeaconCommittee, error) {
	tx := s.tx(ctx)
	if tx == nil {
		ctx, cancel, err := s.BeginTx(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to begin transaction")
		}
		tx = s.tx(ctx)
		defer cancel()
	}

	committee := &chaindb.BeaconCommittee{}

	err := tx.QueryRow(ctx, `
      SELECT f_slot
            ,f_index
            ,f_committee
      FROM t_beacon_committees
      WHERE f_slot = $1
        AND f_index = $2`,
		slot,
		index,
	).Scan(
		&committee.Slot,
		&committee.Index,
		&committee.Committee,
	)
	if err != nil {
		return nil, err
	}
	return committee, nil
}

/*
// AttesterDuties fetches the attester duties at the given slot range for the given validator indices.
func (s *Service) AttesterDuties(ctx context.Context, startSlot spec.Slot, endSlot spec.Slot, validatorIndices []spec.ValidatorIndex) ([]*chaindb.AttesterDuty, error) {
	tx := s.tx(ctx)
	if tx == nil {
		ctx, cancel, err := s.BeginTx(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to begin transaction")
		}
		tx = s.tx(ctx)
		defer cancel()
	}

	rows, err := tx.Query(ctx, `
      SELECT f_slot
            ,f_index
            ,f_committee_index
      FROM t_beacon_committees
      WHERE f_slot >= $1
        AND f_slot < $2
        AND $3 && f_committee
        ORDER BY f_slot, f_index`,
		startSlot,
		endSlot,
		validatorIndices,
	)
	if err != nil {
		return nil, err
	}

	res := make([]*chaindb.AttesterDuty, 0)

	// Map for fast lookup of validator indices.
	validatorIndicesMap := make(map[spec.ValidatorIndex]bool, len(validatorIndices))
	for _, validatorIndex := range validatorIndices {
		validatorIndicesMap[validatorIndex] = true
	}

	for rows.Next() {
		var slot uint64
		var index uint64
		var committee []uint64
		err := rows.Scan(
			&slot,
			&index,
			&committee,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		for i, validatorIndex := range committee {
			if validatorIndicesMap[spec.ValidatorIndex(validatorIndex)] {
				res = append(res, &chaindb.AttesterDuty{
					Slot:           spec.Slot(slot),
					Committee:      spec.CommitteeIndex(index),
					ValidatorIndex: spec.ValidatorIndex(validatorIndex),
					CommitteeIndex: uint64(i),
				})
			}
		}
	}

	return res, nil
}
*/
