package usecase

import (
	"errors"
	"strings"
	"time"

	"asset-management-service/internal/domain"
	"asset-management-service/internal/repository"
	"asset-management-service/pkg/client"

	"gorm.io/gorm"
)

const (
	AssetTypeFolder = "folder"
	AssetTypeNote   = "note"

	AccessRead  = "read"
	AccessWrite = "write"
)

var (
	ErrFolderNameRequired      = errors.New("folder name is required")
	ErrNoteTitleRequired       = errors.New("note title is required")
	ErrInvalidAssetType        = errors.New("asset type must be folder or note")
	ErrInvalidAccessLevel      = errors.New("access level must be read or write")
	ErrFolderNotFound          = errors.New("folder not found")
	ErrNoteNotFound            = errors.New("note not found")
	ErrShareNotFound           = errors.New("share not found")
	ErrForbidden               = errors.New("forbidden")
	ErrCannotShareWithSelf     = errors.New("cannot share asset with yourself")
	ErrShareTargetNotFound     = errors.New("shared target user not found")
	ErrUnableToCheckTeamAccess = errors.New("unable to verify manager oversight")
)

type CreateFolderRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateFolderRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type ShareAssetRequest struct {
	AssetType        string `json:"assetType" binding:"required"`
	AssetID          uint   `json:"assetId" binding:"required"`
	SharedWithUserID uint   `json:"sharedWithUserId" binding:"required"`
	AccessLevel      string `json:"accessLevel" binding:"required"`
}

type FolderResponse struct {
	FolderID    uint      `json:"folderId"`
	OwnerUserID uint      `json:"ownerUserId"`
	Name        string    `json:"name"`
	AccessLevel string    `json:"accessLevel"`
	CanWrite    bool      `json:"canWrite"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type NoteResponse struct {
	NoteID      uint      `json:"noteId"`
	FolderID    uint      `json:"folderId"`
	OwnerUserID uint      `json:"ownerUserId"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	AccessLevel string    `json:"accessLevel"`
	CanWrite    bool      `json:"canWrite"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ShareResponse struct {
	ShareID          uint      `json:"shareId"`
	AssetType        string    `json:"assetType"`
	AssetID          uint      `json:"assetId"`
	SharedByUserID   uint      `json:"sharedByUserId"`
	SharedWithUserID uint      `json:"sharedWithUserId"`
	AccessLevel      string    `json:"accessLevel"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type AssetUsecase interface {
	CreateFolder(actorID uint, req *CreateFolderRequest) (*FolderResponse, error)
	ListFolders(actorID uint, actorRole string, token string) ([]FolderResponse, error)
	GetFolder(actorID uint, actorRole string, token string, folderID uint) (*FolderResponse, error)
	UpdateFolder(actorID uint, actorRole string, token string, folderID uint, req *UpdateFolderRequest) (*FolderResponse, error)
	DeleteFolder(actorID uint, actorRole string, token string, folderID uint) error

	CreateNote(actorID uint, actorRole string, token string, folderID uint, req *CreateNoteRequest) (*NoteResponse, error)
	ListNotesByFolder(actorID uint, actorRole string, token string, folderID uint) ([]NoteResponse, error)
	GetNote(actorID uint, actorRole string, token string, noteID uint) (*NoteResponse, error)
	UpdateNote(actorID uint, actorRole string, token string, noteID uint, req *UpdateNoteRequest) (*NoteResponse, error)
	DeleteNote(actorID uint, actorRole string, token string, noteID uint) error

	ShareAsset(actorID uint, actorRole string, token string, req *ShareAssetRequest) (*ShareResponse, error)
	RevokeShare(actorID uint, actorRole string, token string, shareID uint) error
	ListReceivedShares(actorID uint) ([]ShareResponse, error)
	ListGrantedShares(actorID uint) ([]ShareResponse, error)
}

type assetUsecaseImpl struct {
	repo       repository.AssetRepository
	authClient client.AuthClient
	teamClient client.TeamClient
}

func NewAssetUsecase(repo repository.AssetRepository, authClient client.AuthClient, teamClient client.TeamClient) AssetUsecase {
	return &assetUsecaseImpl{
		repo:       repo,
		authClient: authClient,
		teamClient: teamClient,
	}
}

func (u *assetUsecaseImpl) CreateFolder(actorID uint, req *CreateFolderRequest) (*FolderResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrFolderNameRequired
	}

	folder := &domain.Folder{
		OwnerUserID: actorID,
		Name:        strings.TrimSpace(req.Name),
	}
	if err := u.repo.CreateFolder(folder); err != nil {
		return nil, err
	}

	return toFolderResponse(folder, AccessWrite, true), nil
}

func (u *assetUsecaseImpl) ListFolders(actorID uint, actorRole string, token string) ([]FolderResponse, error) {
	owned, err := u.repo.ListFoldersByOwner(actorID)
	if err != nil {
		return nil, err
	}

	shared, err := u.repo.ListFoldersBySharedUser(actorID)
	if err != nil {
		return nil, err
	}

	folderMap := map[uint]domain.Folder{}
	for i := range owned {
		folderMap[owned[i].ID] = owned[i]
	}
	for i := range shared {
		folderMap[shared[i].ID] = shared[i]
	}

	if actorRole == "manager" && token != "" {
		teams, teamErr := u.teamClient.ListMyTeams(token)
		if teamErr == nil {
			memberIDs := map[uint]struct{}{}
			for i := range teams {
				if !containsUint(teams[i].Managers, actorID) {
					continue
				}
				for j := range teams[i].Members {
					memberIDs[teams[i].Members[j]] = struct{}{}
				}
			}

			for memberID := range memberIDs {
				memberFolders, memberErr := u.repo.ListFoldersByOwner(memberID)
				if memberErr != nil {
					continue
				}
				for i := range memberFolders {
					folderMap[memberFolders[i].ID] = memberFolders[i]
				}
			}
		}
	}

	result := make([]FolderResponse, 0, len(folderMap))
	for _, folder := range folderMap {
		canRead, canWrite, accessErr := u.resolveFolderAccess(actorID, actorRole, token, &folder)
		if accessErr != nil || !canRead {
			continue
		}
		result = append(result, *toFolderResponse(&folder, resolveAccessLevel(canRead, canWrite), canWrite))
	}

	return result, nil
}

func (u *assetUsecaseImpl) GetFolder(actorID uint, actorRole string, token string, folderID uint) (*FolderResponse, error) {
	folder, err := u.repo.FindFolderByID(folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFolderNotFound
		}
		return nil, err
	}

	canRead, canWrite, err := u.resolveFolderAccess(actorID, actorRole, token, folder)
	if err != nil {
		return nil, err
	}
	if !canRead {
		return nil, ErrForbidden
	}

	return toFolderResponse(folder, resolveAccessLevel(canRead, canWrite), canWrite), nil
}

func (u *assetUsecaseImpl) UpdateFolder(actorID uint, actorRole string, token string, folderID uint, req *UpdateFolderRequest) (*FolderResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrFolderNameRequired
	}

	folder, err := u.repo.FindFolderByID(folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFolderNotFound
		}
		return nil, err
	}

	_, canWrite, err := u.resolveFolderAccess(actorID, actorRole, token, folder)
	if err != nil {
		return nil, err
	}
	if !canWrite {
		return nil, ErrForbidden
	}

	folder.Name = strings.TrimSpace(req.Name)
	if err := u.repo.UpdateFolder(folder); err != nil {
		return nil, err
	}
	return toFolderResponse(folder, AccessWrite, true), nil
}

func (u *assetUsecaseImpl) DeleteFolder(actorID uint, actorRole string, token string, folderID uint) error {
	folder, err := u.repo.FindFolderByID(folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFolderNotFound
		}
		return err
	}

	_, canWrite, err := u.resolveFolderAccess(actorID, actorRole, token, folder)
	if err != nil {
		return err
	}
	if !canWrite {
		return ErrForbidden
	}

	return u.repo.DeleteFolder(folderID)
}

func (u *assetUsecaseImpl) CreateNote(actorID uint, actorRole string, token string, folderID uint, req *CreateNoteRequest) (*NoteResponse, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, ErrNoteTitleRequired
	}

	folder, err := u.repo.FindFolderByID(folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFolderNotFound
		}
		return nil, err
	}

	_, canWrite, err := u.resolveFolderAccess(actorID, actorRole, token, folder)
	if err != nil {
		return nil, err
	}
	if !canWrite {
		return nil, ErrForbidden
	}

	note := &domain.Note{
		FolderID:    folderID,
		OwnerUserID: folder.OwnerUserID,
		Title:       strings.TrimSpace(req.Title),
		Content:     req.Content,
	}
	if err := u.repo.CreateNote(note); err != nil {
		return nil, err
	}

	return toNoteResponse(note, AccessWrite, true), nil
}

func (u *assetUsecaseImpl) ListNotesByFolder(actorID uint, actorRole string, token string, folderID uint) ([]NoteResponse, error) {
	folder, err := u.repo.FindFolderByID(folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFolderNotFound
		}
		return nil, err
	}

	canRead, canWrite, err := u.resolveFolderAccess(actorID, actorRole, token, folder)
	if err != nil {
		return nil, err
	}
	if !canRead {
		return nil, ErrForbidden
	}

	notes, err := u.repo.ListNotesByFolder(folderID)
	if err != nil {
		return nil, err
	}

	res := make([]NoteResponse, 0, len(notes))
	for i := range notes {
		res = append(res, *toNoteResponse(&notes[i], resolveAccessLevel(canRead, canWrite), canWrite))
	}
	return res, nil
}

func (u *assetUsecaseImpl) GetNote(actorID uint, actorRole string, token string, noteID uint) (*NoteResponse, error) {
	note, err := u.repo.FindNoteByID(noteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}

	canRead, canWrite, err := u.resolveNoteAccess(actorID, actorRole, token, note)
	if err != nil {
		return nil, err
	}
	if !canRead {
		return nil, ErrForbidden
	}

	return toNoteResponse(note, resolveAccessLevel(canRead, canWrite), canWrite), nil
}

func (u *assetUsecaseImpl) UpdateNote(actorID uint, actorRole string, token string, noteID uint, req *UpdateNoteRequest) (*NoteResponse, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, ErrNoteTitleRequired
	}

	note, err := u.repo.FindNoteByID(noteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}

	_, canWrite, err := u.resolveNoteAccess(actorID, actorRole, token, note)
	if err != nil {
		return nil, err
	}
	if !canWrite {
		return nil, ErrForbidden
	}

	note.Title = strings.TrimSpace(req.Title)
	note.Content = req.Content
	if err := u.repo.UpdateNote(note); err != nil {
		return nil, err
	}

	return toNoteResponse(note, AccessWrite, true), nil
}

func (u *assetUsecaseImpl) DeleteNote(actorID uint, actorRole string, token string, noteID uint) error {
	note, err := u.repo.FindNoteByID(noteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoteNotFound
		}
		return err
	}

	_, canWrite, err := u.resolveNoteAccess(actorID, actorRole, token, note)
	if err != nil {
		return err
	}
	if !canWrite {
		return ErrForbidden
	}

	return u.repo.DeleteNote(noteID)
}

func (u *assetUsecaseImpl) ShareAsset(actorID uint, actorRole string, token string, req *ShareAssetRequest) (*ShareResponse, error) {
	assetType := strings.TrimSpace(strings.ToLower(req.AssetType))
	accessLevel := strings.TrimSpace(strings.ToLower(req.AccessLevel))

	if assetType != AssetTypeFolder && assetType != AssetTypeNote {
		return nil, ErrInvalidAssetType
	}
	if accessLevel != AccessRead && accessLevel != AccessWrite {
		return nil, ErrInvalidAccessLevel
	}
	if req.SharedWithUserID == actorID {
		return nil, ErrCannotShareWithSelf
	}

	users, err := u.authClient.GetUsers(token, []uint{req.SharedWithUserID})
	if err != nil || len(users) == 0 {
		return nil, ErrShareTargetNotFound
	}

	switch assetType {
	case AssetTypeFolder:
		folder, folderErr := u.repo.FindFolderByID(req.AssetID)
		if folderErr != nil {
			if errors.Is(folderErr, gorm.ErrRecordNotFound) {
				return nil, ErrFolderNotFound
			}
			return nil, folderErr
		}
		_, canWrite, accessErr := u.resolveFolderAccess(actorID, actorRole, token, folder)
		if accessErr != nil {
			return nil, accessErr
		}
		if !canWrite {
			return nil, ErrForbidden
		}
	case AssetTypeNote:
		note, noteErr := u.repo.FindNoteByID(req.AssetID)
		if noteErr != nil {
			if errors.Is(noteErr, gorm.ErrRecordNotFound) {
				return nil, ErrNoteNotFound
			}
			return nil, noteErr
		}
		_, canWrite, accessErr := u.resolveNoteAccess(actorID, actorRole, token, note)
		if accessErr != nil {
			return nil, accessErr
		}
		if !canWrite {
			return nil, ErrForbidden
		}
	}

	share := &domain.AssetShare{
		AssetType:        assetType,
		AssetID:          req.AssetID,
		SharedByUserID:   actorID,
		SharedWithUserID: req.SharedWithUserID,
		AccessLevel:      accessLevel,
	}
	if err := u.repo.UpsertShare(share); err != nil {
		return nil, err
	}

	currentShare, err := u.repo.FindShare(assetType, req.AssetID, req.SharedWithUserID)
	if err != nil {
		return nil, err
	}
	return toShareResponse(currentShare), nil
}

func (u *assetUsecaseImpl) RevokeShare(actorID uint, actorRole string, token string, shareID uint) error {
	share, err := u.repo.FindShareByID(shareID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrShareNotFound
		}
		return err
	}

	if share.SharedByUserID == actorID {
		return u.repo.DeleteShareByID(shareID)
	}

	switch share.AssetType {
	case AssetTypeFolder:
		folder, folderErr := u.repo.FindFolderByID(share.AssetID)
		if folderErr != nil {
			return ErrFolderNotFound
		}
		_, canWrite, accessErr := u.resolveFolderAccess(actorID, actorRole, token, folder)
		if accessErr != nil {
			return accessErr
		}
		if !canWrite {
			return ErrForbidden
		}
	case AssetTypeNote:
		note, noteErr := u.repo.FindNoteByID(share.AssetID)
		if noteErr != nil {
			return ErrNoteNotFound
		}
		_, canWrite, accessErr := u.resolveNoteAccess(actorID, actorRole, token, note)
		if accessErr != nil {
			return accessErr
		}
		if !canWrite {
			return ErrForbidden
		}
	default:
		return ErrInvalidAssetType
	}

	return u.repo.DeleteShareByID(shareID)
}

func (u *assetUsecaseImpl) ListReceivedShares(actorID uint) ([]ShareResponse, error) {
	shares, err := u.repo.ListSharesReceived(actorID)
	if err != nil {
		return nil, err
	}

	result := make([]ShareResponse, 0, len(shares))
	for i := range shares {
		result = append(result, *toShareResponse(&shares[i]))
	}
	return result, nil
}

func (u *assetUsecaseImpl) ListGrantedShares(actorID uint) ([]ShareResponse, error) {
	shares, err := u.repo.ListSharesGranted(actorID)
	if err != nil {
		return nil, err
	}

	result := make([]ShareResponse, 0, len(shares))
	for i := range shares {
		result = append(result, *toShareResponse(&shares[i]))
	}
	return result, nil
}

func (u *assetUsecaseImpl) resolveFolderAccess(actorID uint, actorRole string, token string, folder *domain.Folder) (bool, bool, error) {
	if folder.OwnerUserID == actorID {
		return true, true, nil
	}

	share, err := u.repo.FindShare(AssetTypeFolder, folder.ID, actorID)
	if err == nil && share != nil {
		if share.AccessLevel == AccessWrite {
			return true, true, nil
		}
		return true, false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, false, err
	}

	isOversight, oversightErr := u.isManagerOversight(actorID, actorRole, token, folder.OwnerUserID)
	if oversightErr != nil {
		return false, false, oversightErr
	}
	if isOversight {
		return true, false, nil
	}

	return false, false, nil
}

func (u *assetUsecaseImpl) resolveNoteAccess(actorID uint, actorRole string, token string, note *domain.Note) (bool, bool, error) {
	if note.OwnerUserID == actorID {
		return true, true, nil
	}

	canRead := false
	canWrite := false

	share, err := u.repo.FindShare(AssetTypeNote, note.ID, actorID)
	if err == nil && share != nil {
		canRead = true
		if share.AccessLevel == AccessWrite {
			canWrite = true
		}
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, false, err
	}

	folder, folderErr := u.repo.FindFolderByID(note.FolderID)
	if folderErr != nil {
		if errors.Is(folderErr, gorm.ErrRecordNotFound) {
			return false, false, ErrFolderNotFound
		}
		return false, false, folderErr
	}
	folderRead, folderWrite, accessErr := u.resolveFolderAccess(actorID, actorRole, token, folder)
	if accessErr != nil {
		return false, false, accessErr
	}

	if folderRead {
		canRead = true
	}
	if folderWrite {
		canWrite = true
	}

	return canRead, canWrite, nil
}

func (u *assetUsecaseImpl) isManagerOversight(actorID uint, actorRole string, token string, ownerUserID uint) (bool, error) {
	if actorRole != "manager" || actorID == ownerUserID {
		return false, nil
	}
	if token == "" {
		return false, ErrUnableToCheckTeamAccess
	}

	isManager, err := u.teamClient.IsManagerOf(token, actorID, ownerUserID)
	if err != nil {
		return false, ErrUnableToCheckTeamAccess
	}
	return isManager, nil
}

func toFolderResponse(folder *domain.Folder, accessLevel string, canWrite bool) *FolderResponse {
	return &FolderResponse{
		FolderID:    folder.ID,
		OwnerUserID: folder.OwnerUserID,
		Name:        folder.Name,
		AccessLevel: accessLevel,
		CanWrite:    canWrite,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}
}

func toNoteResponse(note *domain.Note, accessLevel string, canWrite bool) *NoteResponse {
	return &NoteResponse{
		NoteID:      note.ID,
		FolderID:    note.FolderID,
		OwnerUserID: note.OwnerUserID,
		Title:       note.Title,
		Content:     note.Content,
		AccessLevel: accessLevel,
		CanWrite:    canWrite,
		CreatedAt:   note.CreatedAt,
		UpdatedAt:   note.UpdatedAt,
	}
}

func toShareResponse(share *domain.AssetShare) *ShareResponse {
	return &ShareResponse{
		ShareID:          share.ID,
		AssetType:        share.AssetType,
		AssetID:          share.AssetID,
		SharedByUserID:   share.SharedByUserID,
		SharedWithUserID: share.SharedWithUserID,
		AccessLevel:      share.AccessLevel,
		CreatedAt:        share.CreatedAt,
		UpdatedAt:        share.UpdatedAt,
	}
}

func resolveAccessLevel(canRead bool, canWrite bool) string {
	if canWrite {
		return AccessWrite
	}
	if canRead {
		return AccessRead
	}
	return ""
}

func containsUint(values []uint, target uint) bool {
	for i := range values {
		if values[i] == target {
			return true
		}
	}
	return false
}
