// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     CoreRegistryJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package schemas

import (
	"os"
	"path/filepath"
	"runtime"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

type CoreKind struct {
	Name    string
	CueFile cue.Value
}

func GetCoreKinds() ([]CoreKind, error) {
	ctx := cuecontext.New()
	kinds := make([]CoreKind, 0)

	_, caller, _, _ := runtime.Caller(0)
	root := filepath.Join(caller, "../../../..")

	accesspolicyCue, err := loadCueFile(ctx, filepath.Join(root, "./kinds/accesspolicy/access_policy_kind.cue"))
	if err != nil {
		return nil, err
	}
	kinds = append(kinds, CoreKind{
		Name:    "accesspolicy",
		CueFile: accesspolicyCue,
	})

	dashboardCue, err := loadCueFile(ctx, filepath.Join(root, "./kinds/dashboard/dashboard_kind.cue"))
	if err != nil {
		return nil, err
	}
	kinds = append(kinds, CoreKind{
		Name:    "dashboard",
		CueFile: dashboardCue,
	})

	librarypanelCue, err := loadCueFile(ctx, filepath.Join(root, "./kinds/librarypanel/librarypanel_kind.cue"))
	if err != nil {
		return nil, err
	}
	kinds = append(kinds, CoreKind{
		Name:    "librarypanel",
		CueFile: librarypanelCue,
	})

	preferencesCue, err := loadCueFile(ctx, filepath.Join(root, "./kinds/preferences/preferences_kind.cue"))
	if err != nil {
		return nil, err
	}
	kinds = append(kinds, CoreKind{
		Name:    "preferences",
		CueFile: preferencesCue,
	})

	publicdashboardCue, err := loadCueFile(ctx, filepath.Join(root, "./kinds/publicdashboard/public_dashboard_kind.cue"))
	if err != nil {
		return nil, err
	}
	kinds = append(kinds, CoreKind{
		Name:    "publicdashboard",
		CueFile: publicdashboardCue,
	})

	roleCue, err := loadCueFile(ctx, filepath.Join(root, "./kinds/role/role_kind.cue"))
	if err != nil {
		return nil, err
	}
	kinds = append(kinds, CoreKind{
		Name:    "role",
		CueFile: roleCue,
	})

	rolebindingCue, err := loadCueFile(ctx, filepath.Join(root, "./kinds/rolebinding/role_binding_kind.cue"))
	if err != nil {
		return nil, err
	}
	kinds = append(kinds, CoreKind{
		Name:    "rolebinding",
		CueFile: rolebindingCue,
	})

	return kinds, nil
}

func loadCueFile(ctx *cue.Context, path string) (cue.Value, error) {
	cueFile, err := os.ReadFile(path)
	if err != nil {
		return cue.Value{}, err
	}

	return ctx.CompileBytes(cueFile), nil
}
