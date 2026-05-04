package repository

import (
	"asset-management-service/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AssetRepository interface {
	CreateFolder(folder *domain.Folder) error
	FindFolderByID(folderID uint) (*domain.Folder, error)
	ListFoldersByOwner(ownerUserID uint) ([]domain.Folder, error)
	ListFoldersBySharedUser(userID uint) ([]domain.Folder, error)
	UpdateFolder(folder *domain.Folder) error
	DeleteFolder(folderID uint) error

	CreateNote(note *domain.Note) error
	FindNoteByID(noteID uint) (*domain.Note, error)
	ListNotesByFolder(folderID uint) ([]domain.Note, error)
	UpdateNote(note *domain.Note) error
	DeleteNote(noteID uint) error

	FindShare(assetType string, assetID uint, sharedWithUserID uint) (*domain.AssetShare, error)
	UpsertShare(share *domain.AssetShare) error
	DeleteShareByID(shareID uint) error
	FindShareByID(shareID uint) (*domain.AssetShare, error)
	ListSharesReceived(userID uint) ([]domain.AssetShare, error)
	ListSharesGranted(userID uint) ([]domain.AssetShare, error)
}

type assetRepositoryImpl struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepositoryImpl{db: db}
}

func (r *assetRepositoryImpl) CreateFolder(folder *domain.Folder) error {
	return r.db.Create(folder).Error
}

func (r *assetRepositoryImpl) FindFolderByID(folderID uint) (*domain.Folder, error) {
	var folder domain.Folder
	err := r.db.First(&folder, folderID).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *assetRepositoryImpl) ListFoldersByOwner(ownerUserID uint) ([]domain.Folder, error) {
	var folders []domain.Folder
	err := r.db.Where("owner_user_id = ?", ownerUserID).Order("id DESC").Find(&folders).Error
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (r *assetRepositoryImpl) ListFoldersBySharedUser(userID uint) ([]domain.Folder, error) {
	var folders []domain.Folder
	err := r.db.Model(&domain.Folder{}).
		Joins("JOIN asset_shares ON asset_shares.asset_type = ? AND asset_shares.asset_id = folders.id", "folder").
		Where("asset_shares.shared_with_user_id = ?", userID).
		Order("folders.id DESC").
		Find(&folders).Error
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (r *assetRepositoryImpl) UpdateFolder(folder *domain.Folder) error {
	return r.db.Save(folder).Error
}

func (r *assetRepositoryImpl) DeleteFolder(folderID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var notes []domain.Note
		if err := tx.Where("folder_id = ?", folderID).Find(&notes).Error; err != nil {
			return err
		}

		if err := tx.Where("asset_type = ? AND asset_id = ?", "folder", folderID).Delete(&domain.AssetShare{}).Error; err != nil {
			return err
		}

		if len(notes) > 0 {
			noteIDs := make([]uint, 0, len(notes))
			for i := range notes {
				noteIDs = append(noteIDs, notes[i].ID)
			}
			if err := tx.Where("asset_type = ? AND asset_id IN ?", "note", noteIDs).Delete(&domain.AssetShare{}).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("folder_id = ?", folderID).Delete(&domain.Note{}).Error; err != nil {
			return err
		}

		return tx.Delete(&domain.Folder{}, folderID).Error
	})
}

func (r *assetRepositoryImpl) CreateNote(note *domain.Note) error {
	return r.db.Create(note).Error
}

func (r *assetRepositoryImpl) FindNoteByID(noteID uint) (*domain.Note, error) {
	var note domain.Note
	err := r.db.First(&note, noteID).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *assetRepositoryImpl) ListNotesByFolder(folderID uint) ([]domain.Note, error) {
	var notes []domain.Note
	err := r.db.Where("folder_id = ?", folderID).Order("id DESC").Find(&notes).Error
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *assetRepositoryImpl) UpdateNote(note *domain.Note) error {
	return r.db.Save(note).Error
}

func (r *assetRepositoryImpl) DeleteNote(noteID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("asset_type = ? AND asset_id = ?", "note", noteID).Delete(&domain.AssetShare{}).Error; err != nil {
			return err
		}
		return tx.Delete(&domain.Note{}, noteID).Error
	})
}

func (r *assetRepositoryImpl) FindShare(assetType string, assetID uint, sharedWithUserID uint) (*domain.AssetShare, error) {
	var share domain.AssetShare
	err := r.db.Where("asset_type = ? AND asset_id = ? AND shared_with_user_id = ?", assetType, assetID, sharedWithUserID).
		First(&share).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *assetRepositoryImpl) UpsertShare(share *domain.AssetShare) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "asset_type"},
			{Name: "asset_id"},
			{Name: "shared_with_user_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"shared_by_user_id", "access_level", "updated_at"}),
	}).Create(share).Error
}

func (r *assetRepositoryImpl) DeleteShareByID(shareID uint) error {
	return r.db.Delete(&domain.AssetShare{}, shareID).Error
}

func (r *assetRepositoryImpl) FindShareByID(shareID uint) (*domain.AssetShare, error) {
	var share domain.AssetShare
	if err := r.db.First(&share, shareID).Error; err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *assetRepositoryImpl) ListSharesReceived(userID uint) ([]domain.AssetShare, error) {
	var shares []domain.AssetShare
	err := r.db.Where("shared_with_user_id = ?", userID).Order("id DESC").Find(&shares).Error
	if err != nil {
		return nil, err
	}
	return shares, nil
}

func (r *assetRepositoryImpl) ListSharesGranted(userID uint) ([]domain.AssetShare, error) {
	var shares []domain.AssetShare
	err := r.db.Where("shared_by_user_id = ?", userID).Order("id DESC").Find(&shares).Error
	if err != nil {
		return nil, err
	}
	return shares, nil
}
