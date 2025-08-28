package workspace

import (
	"altron/core/generated/spec"
	"altron/core/repositories/ent"
)

func PresentWorkspace(workspaceEnt *ent.Workspace) *spec.Workspace {
	return &spec.Workspace{
		Id:     workspaceEnt.ID.String(),
		Name:   workspaceEnt.Name,
		Status: spec.WorkspaceStatus(workspaceEnt.Status),
	}
}
