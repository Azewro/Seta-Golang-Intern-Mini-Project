package usecase

import (
	"errors"
	"strings"
	"time"

	"team-management-service/internal/domain"
	"team-management-service/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrTeamNameRequired      = errors.New("team name is required")
	ErrForbidden             = errors.New("forbidden")
	ErrTeamNotFound          = errors.New("team not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrAlreadyInTeam         = errors.New("user already belongs to the team")
	ErrNotTeamMember         = errors.New("user does not belong to this team")
	ErrNotTeamManager        = errors.New("team manager role required")
	ErrMainManagerRequired   = errors.New("main manager role required")
	ErrGlobalManagerRequired = errors.New("target user must have global manager role")
	ErrCannotRemoveMain      = errors.New("cannot remove main manager from managers list")
	ErrUseManagerEndpoint    = errors.New("target user is a manager, use manager endpoint")
)

// CreateTeamRequest defines payload for team creation.
type CreateTeamRequest struct {
	TeamName string `json:"teamName" binding:"required"`
}

// TeamActionRequest defines payload for user assignment actions.
type TeamActionRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

// TeamResponse defines team output payload.
type TeamResponse struct {
	TeamID            uint      `json:"teamId"`
	TeamName          string    `json:"teamName"`
	MainManagerUserID uint      `json:"mainManagerUserId"`
	Managers          []uint    `json:"managers"`
	Members           []uint    `json:"members"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// TeamUsecase exposes team business logic.
type TeamUsecase interface {
	CreateTeam(actorID uint, actorRole string, req *CreateTeamRequest) (*TeamResponse, error)
	GetTeam(actorID uint, teamID uint) (*TeamResponse, error)
	ListMyTeams(actorID uint) ([]TeamResponse, error)
	AddMember(actorID uint, actorRole string, teamID uint, targetUserID uint) error
	RemoveMember(actorID uint, actorRole string, teamID uint, targetUserID uint) error
	AddManager(actorID uint, actorRole string, teamID uint, targetUserID uint) error
	RemoveManager(actorID uint, actorRole string, teamID uint, targetUserID uint) error
}

type teamUsecaseImpl struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewTeamUsecase(teamRepo repository.TeamRepository, userRepo repository.UserRepository) TeamUsecase {
	return &teamUsecaseImpl{teamRepo: teamRepo, userRepo: userRepo}
}

func (u *teamUsecaseImpl) CreateTeam(actorID uint, actorRole string, req *CreateTeamRequest) (*TeamResponse, error) {
	if actorRole != "manager" {
		return nil, ErrForbidden
	}
	if strings.TrimSpace(req.TeamName) == "" {
		return nil, ErrTeamNameRequired
	}

	team := &domain.Team{
		TeamName:          strings.TrimSpace(req.TeamName),
		MainManagerUserID: actorID,
	}
	if err := u.teamRepo.CreateTeamWithMainManager(team); err != nil {
		return nil, err
	}

	return u.buildTeamResponse(team)
}

func (u *teamUsecaseImpl) GetTeam(actorID uint, teamID uint) (*TeamResponse, error) {
	team, err := u.teamRepo.FindTeamByID(teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}

	if err := u.requireInTeam(teamID, actorID); err != nil {
		return nil, err
	}

	return u.buildTeamResponse(team)
}

func (u *teamUsecaseImpl) ListMyTeams(actorID uint) ([]TeamResponse, error) {
	teams, err := u.teamRepo.ListTeamsByUser(actorID)
	if err != nil {
		return nil, err
	}

	response := make([]TeamResponse, 0, len(teams))
	for i := range teams {
		teamResponse, buildErr := u.buildTeamResponse(&teams[i])
		if buildErr != nil {
			return nil, buildErr
		}
		response = append(response, *teamResponse)
	}
	return response, nil
}

func (u *teamUsecaseImpl) AddMember(actorID uint, actorRole string, teamID uint, targetUserID uint) error {
	if actorRole != "manager" {
		return ErrForbidden
	}
	if err := u.requireTeamManager(teamID, actorID); err != nil {
		return err
	}

	if _, err := u.userRepo.FindByID(targetUserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	membership, err := u.teamRepo.FindMembership(teamID, targetUserID)
	if err == nil && membership != nil {
		return ErrAlreadyInTeam
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return u.teamRepo.CreateMembership(teamID, targetUserID, "member")
}

func (u *teamUsecaseImpl) RemoveMember(actorID uint, actorRole string, teamID uint, targetUserID uint) error {
	if actorRole != "manager" {
		return ErrForbidden
	}
	if err := u.requireTeamManager(teamID, actorID); err != nil {
		return err
	}

	membership, err := u.teamRepo.FindMembership(teamID, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotTeamMember
		}
		return err
	}
	if membership.MembershipRole == "manager" {
		return ErrUseManagerEndpoint
	}

	return u.teamRepo.DeleteMembership(teamID, targetUserID)
}

func (u *teamUsecaseImpl) AddManager(actorID uint, actorRole string, teamID uint, targetUserID uint) error {
	if actorRole != "manager" {
		return ErrForbidden
	}
	if err := u.requireMainManager(teamID, actorID); err != nil {
		return err
	}

	targetUser, err := u.userRepo.FindByID(targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	if targetUser.Role != "manager" {
		return ErrGlobalManagerRequired
	}

	membership, err := u.teamRepo.FindMembership(teamID, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u.teamRepo.CreateMembership(teamID, targetUserID, "manager")
		}
		return err
	}

	if membership.MembershipRole == "manager" {
		return ErrAlreadyInTeam
	}

	return u.teamRepo.UpdateMembershipRole(teamID, targetUserID, "manager")
}

func (u *teamUsecaseImpl) RemoveManager(actorID uint, actorRole string, teamID uint, targetUserID uint) error {
	if actorRole != "manager" {
		return ErrForbidden
	}
	if err := u.requireMainManager(teamID, actorID); err != nil {
		return err
	}

	team, err := u.teamRepo.FindTeamByID(teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTeamNotFound
		}
		return err
	}
	if team.MainManagerUserID == targetUserID {
		return ErrCannotRemoveMain
	}

	membership, err := u.teamRepo.FindMembership(teamID, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotTeamMember
		}
		return err
	}
	if membership.MembershipRole != "manager" {
		return ErrNotTeamManager
	}

	return u.teamRepo.DeleteMembership(teamID, targetUserID)
}

func (u *teamUsecaseImpl) buildTeamResponse(team *domain.Team) (*TeamResponse, error) {
	memberships, err := u.teamRepo.ListMembershipsByTeam(team.ID)
	if err != nil {
		return nil, err
	}

	managers := make([]uint, 0)
	members := make([]uint, 0)
	for i := range memberships {
		if memberships[i].MembershipRole == "manager" {
			managers = append(managers, memberships[i].UserID)
			continue
		}
		members = append(members, memberships[i].UserID)
	}

	return &TeamResponse{
		TeamID:            team.ID,
		TeamName:          team.TeamName,
		MainManagerUserID: team.MainManagerUserID,
		Managers:          managers,
		Members:           members,
		CreatedAt:         team.CreatedAt,
		UpdatedAt:         team.UpdatedAt,
	}, nil
}

func (u *teamUsecaseImpl) requireInTeam(teamID uint, userID uint) error {
	_, err := u.teamRepo.FindMembership(teamID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrForbidden
		}
		return err
	}
	return nil
}

func (u *teamUsecaseImpl) requireTeamManager(teamID uint, userID uint) error {
	membership, err := u.teamRepo.FindMembership(teamID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrForbidden
		}
		return err
	}
	if membership.MembershipRole != "manager" {
		return ErrForbidden
	}
	return nil
}

func (u *teamUsecaseImpl) requireMainManager(teamID uint, userID uint) error {
	team, err := u.teamRepo.FindTeamByID(teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTeamNotFound
		}
		return err
	}
	if team.MainManagerUserID != userID {
		return ErrMainManagerRequired
	}
	return nil
}
