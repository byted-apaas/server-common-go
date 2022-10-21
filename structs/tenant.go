// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package structs

type Tenant struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Type      int64  `json:"type"`
	Namespace string `json:"namespace"`
	Domain    string `json:"domain"`
}

type Initiator struct {
	ID int64 `json:"_id"`
}
