package lty_models

import (
	"time"

	"lty_backend/lty_config"
)

type ProjectGroup struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateProjectGroup(group *ProjectGroup) error {
	now := time.Now()
	if group.CreatedAt.IsZero() {
		group.CreatedAt = now
	}
	if group.UpdatedAt.IsZero() {
		group.UpdatedAt = now
	}
	query := "INSERT INTO lty_project_groups (user_id, name, created_at, updated_at) VALUES (?, ?, ?, ?)"
	result, err := lty_config.DB.Exec(query, group.UserID, group.Name, group.CreatedAt, group.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	group.ID = uint(id)
	return nil
}

func UpdateProjectGroup(userID uint, groupID uint, name string) error {
	query := "UPDATE lty_project_groups SET name = ? WHERE id = ? AND user_id = ?"
	_, err := lty_config.DB.Exec(query, name, groupID, userID)
	return err
}

func DeleteProjectGroup(userID uint, groupID uint) error {
	// First reset project_id to 0 for all tasks in this group
	resetQuery := "UPDATE lty_tasks SET project_id = 0 WHERE project_id = ? AND user_id = ?"
	_, err := lty_config.DB.Exec(resetQuery, groupID, userID)
	if err != nil {
		return err
	}

	// Then delete the group
	delQuery := "DELETE FROM lty_project_groups WHERE id = ? AND user_id = ?"
	_, err = lty_config.DB.Exec(delQuery, groupID, userID)
	return err
}

func GetProjectGroupsPaginated(userID uint, page, pageSize int) ([]ProjectGroup, int, error) {
	var total int
	err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_project_groups WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := "SELECT id, user_id, name, created_at, updated_at FROM lty_project_groups WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := lty_config.DB.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var groups []ProjectGroup
	for rows.Next() {
		var g ProjectGroup
		if err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, 0, err
		}
		groups = append(groups, g)
	}
	return groups, total, nil
}

func GetProjectGroups(userID uint) ([]ProjectGroup, error) {
	query := "SELECT id, user_id, name, created_at, updated_at FROM lty_project_groups WHERE user_id = ? ORDER BY created_at DESC"
	rows, err := lty_config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []ProjectGroup
	for rows.Next() {
		var g ProjectGroup
		if err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func GetProjectGroupByID(userID uint, projectID int) (*ProjectGroup, error) {
	query := "SELECT id, user_id, name, created_at, updated_at FROM lty_project_groups WHERE user_id = ? AND id = ?"
	var g ProjectGroup
	err := lty_config.DB.QueryRow(query, userID, projectID).Scan(&g.ID, &g.UserID, &g.Name, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func GetGeneratingProjectGroupIDs(userID uint) ([]uint, error) {
	// 获取当前用户下有正在生成任务(status=0)的项目组ID集合
	query := "SELECT DISTINCT project_id FROM lty_tasks WHERE user_id = ? AND status = 0 AND project_id > 0"
	rows, err := lty_config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupIDs []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		groupIDs = append(groupIDs, id)
	}
	return groupIDs, nil
}
