/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package neo3_state_manager

import (
	"fmt"

	"github.com/polynetwork/poly/common"
)

func SerializeStringArray(data []string) []byte {
	sink := common.NewZeroCopySink(nil)
	// serialize
	sink.WriteVarUint(uint64(len(data)))
	for _, v := range data {
		sink.WriteString(v)
	}
	return sink.Bytes()
}

func DeserializeStringArray(data []byte) ([]string, error) {
	if len(data) == 0 {
		return []string{}, nil
	}
	source := common.NewZeroCopySource(data)
	n, eof := source.NextVarUint()
	if eof {
		return nil, fmt.Errorf("source.NextVarUint error")
	}
	result := make([]string, 0, n)
	for i := 0; uint64(i) < n; i++ {
		ss, eof := source.NextString()
		if eof {
			return nil, fmt.Errorf("source.NextString error")
		}
		result = append(result, ss)
	}
	return result, nil
}
