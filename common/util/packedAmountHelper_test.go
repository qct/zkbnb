/*
 * Copyright © 2021 Zkbas Protocol
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package util

import (
	"fmt"
	"math/big"
	"testing"
)

func TestToPackedAmount(t *testing.T) {
	a, _ := new(big.Int).SetString("34359738361", 10)
	amount, err := ToPackedAmount(a)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(amount)
}

func TestToPackedFee(t *testing.T) {
	amount, _ := new(big.Int).SetString("100000000000000", 10)
	fee, err := ToPackedFee(amount)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fee)
}
