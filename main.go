// Copyright (c) 2021 PointGoal
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/pointgoal/workstation/pkg/controller"
	"github.com/pointgoal/workstation/pkg/datastore"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-entry/entry"
)

// This must be declared in order to register registration function into rk context
// otherwise, rk-boot won't able to bootstrap entry automatically from boot config file
func init() {
	rkentry.RegisterEntryRegFunc(datastore.RegisterDataStoreFromConfig)
	rkentry.RegisterEntryRegFunc(controller.RegisterControllerFromConfig)
}

// @title Workstation
// @version 1.0
// @description This is workstation backend with rk-boot.

// @contact.name PointGoal team
// @contact.url https://github.com/pointgoal/workstation
// @contact.email lark@pointgoal.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// Application entrance.
func main() {
	// Create a new boot instance.
	boot := rkboot.NewBoot()

	// Bootstrap
	boot.Bootstrap(context.Background())

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.Background())
}
