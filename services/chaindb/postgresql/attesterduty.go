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

	"github.com/wealdtech/chaind/services/chaindb"
)

// SetAttesterDuty sets a attester duty.
func (s *Service) SetAttesterDuty(ctx context.Context, attesterDuty *chaindb.AttesterDuty) error {
	tx := s.tx(ctx)
	if tx == nil {
		return ErrNoTransaction
	}

	_, err := tx.Exec(ctx, `
      INSERT INTO t_attester_duties(f_slot
                                   ,f_validator_index
                                   ,f_committee_index
                                   ,f_validator_committee_index
                                   ,f_committee_length)
      VALUES($1,$2,$3,$4,$5)
      ON CONFLICT (f_slot, f_validator_index) DO
      UPDATE
      SET f_committee_index = excluded.f_committee_index,
          f_validator_committee_index = excluded.f_validator_committee_index,
          f_committee_length = excluded.f_committee_length
		 `,
		attesterDuty.Slot,
		attesterDuty.ValidatorIndex,
		attesterDuty.CommitteeIndex,
		attesterDuty.ValidatorCommitteeIndex,
		attesterDuty.CommitteeLength,
	)

	return err
}
