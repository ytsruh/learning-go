//go:build mem
// +build mem

package mem

import (
	"ytsruh.com/saas/data/model"
)

type Admin struct {
	requests []model.APIRequest
}

func (a *Admin) LogRequest(reqs []model.APIRequest) error {
	a.requests = append(a.requests, reqs...)
	return nil
}

func (a *Admin) RefreshSession(conn *bool, dbName string) {
}

func (a *Admin) Close() {
}
