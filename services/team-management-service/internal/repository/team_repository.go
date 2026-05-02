package repository

import (
	"errors"

	"team-management-service/internal/domain"

	"gorm.io/gorm"
)

// TeamRepository defines persistence operations for teams and memberships.
type TeamRepository interface {
	CreateTeamWithMainManager(team *domain.Team) error
	FindTeamByID(teamID uint) (*domain.Team, error)
	ListTeamsByUser(userID uint) ([]domain.Team, error)
	ListMembershipsByTeam(teamID uint) ([]domain.TeamMembership, error)
	FindMembership(teamID uint, userID uint) (*domain.TeamMembership, error)
	CreateMembership(teamID uint, userID uint, role string) error
	UpdateMembershipRole(teamID uint, userID uint, role string) error
	DeleteMembership(teamID uint, userID uint) error
}

type teamRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepositoryImpl{db: db}
}

func (r *teamRepositoryImpl) CreateTeamWithMainManager(team *domain.Team) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			return err
		}

		membership := domain.TeamMembership{
			TeamID:         team.ID,
			UserID:         team.MainManagerUserID,
			MembershipRole: "manager",
		}
		return tx.Create(&membership).Error
	})
}

func (r *teamRepositoryImpl) FindTeamByID(teamID uint) (*domain.Team, error) {
	var team domain.Team
	err := r.db.First(&team, teamID).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepositoryImpl) ListTeamsByUser(userID uint) ([]domain.Team, error) {
	var teams []domain.Team
	err := r.db.Table("teams").
		Joins("JOIN team_memberships ON team_memberships.team_id = teams.id").
		Where("team_memberships.user_id = ?", userID).
		Order("teams.id ASC").
		Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *teamRepositoryImpl) ListMembershipsByTeam(teamID uint) ([]domain.TeamMembership, error) {
	var memberships []domain.TeamMembership
	err := r.db.Where("team_id = ?", teamID).Order("user_id ASC").Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *teamRepositoryImpl) FindMembership(teamID uint, userID uint) (*domain.TeamMembership, error) {
	var membership domain.TeamMembership
	err := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&membership).Error
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

func (r *teamRepositoryImpl) CreateMembership(teamID uint, userID uint, role string) error {
	membership := domain.TeamMembership{
		TeamID:         teamID,
		UserID:         userID,
		MembershipRole: role,
	}
	return r.db.Create(&membership).Error
}

func (r *teamRepositoryImpl) UpdateMembershipRole(teamID uint, userID uint, role string) error {
	result := r.db.Model(&domain.TeamMembership{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Update("membership_role", role)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *teamRepositoryImpl) DeleteMembership(teamID uint, userID uint) error {
	result := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&domain.TeamMembership{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("membership not found")
	}
	return nil
}
