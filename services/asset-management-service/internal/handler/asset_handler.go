package handler

import (
	"errors"
	"net/http"
	"strconv"

	"asset-management-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	assetUsecase usecase.AssetUsecase
}

func NewAssetHandler(assetUsecase usecase.AssetUsecase) *AssetHandler {
	return &AssetHandler{assetUsecase: assetUsecase}
}

// CreateFolder godoc
// @Summary Create folder
// @Tags folders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body usecase.CreateFolderRequest true "Folder name"
// @Success 201 {object} usecase.FolderResponse
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/folders [post]
func (h *AssetHandler) CreateFolder(c *gin.Context) {
	userID, _, _, ok := getActorContext(c)
	if !ok {
		return
	}

	var req usecase.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := h.assetUsecase.CreateFolder(userID, &req)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusCreated, folder)
}

// ListFolders godoc
// @Summary List folders
// @Tags folders
// @Security BearerAuth
// @Produce json
// @Success 200 {object} foldersListBody
// @Failure 401 {object} errBody
// @Failure 403 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/folders [get]
func (h *AssetHandler) ListFolders(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	folders, err := h.assetUsecase.ListFolders(userID, role, token)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": folders})
}

// GetFolder godoc
// @Summary Get folder
// @Tags folders
// @Security BearerAuth
// @Produce json
// @Param folderId path int true "Folder ID"
// @Success 200 {object} usecase.FolderResponse
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/folders/{folderId} [get]
func (h *AssetHandler) GetFolder(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	folderID, err := parseUintParam(c, "folderId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, usecaseErr := h.assetUsecase.GetFolder(userID, role, token, folderID)
	if usecaseErr != nil {
		h.handleUsecaseError(c, usecaseErr)
		return
	}
	c.JSON(http.StatusOK, folder)
}

// UpdateFolder godoc
// @Summary Update folder
// @Tags folders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param folderId path int true "Folder ID"
// @Param body body usecase.UpdateFolderRequest true "New name"
// @Success 200 {object} usecase.FolderResponse
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/folders/{folderId} [patch]
func (h *AssetHandler) UpdateFolder(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	folderID, err := parseUintParam(c, "folderId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req usecase.UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, usecaseErr := h.assetUsecase.UpdateFolder(userID, role, token, folderID, &req)
	if usecaseErr != nil {
		h.handleUsecaseError(c, usecaseErr)
		return
	}
	c.JSON(http.StatusOK, folder)
}

// DeleteFolder godoc
// @Summary Delete folder
// @Tags folders
// @Security BearerAuth
// @Produce json
// @Param folderId path int true "Folder ID"
// @Success 200 {object} msgBody
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/folders/{folderId} [delete]
func (h *AssetHandler) DeleteFolder(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	folderID, err := parseUintParam(c, "folderId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if usecaseErr := h.assetUsecase.DeleteFolder(userID, role, token, folderID); usecaseErr != nil {
		h.handleUsecaseError(c, usecaseErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted"})
}

// CreateNote godoc
// @Summary Create note
// @Tags notes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param folderId path int true "Folder ID"
// @Param body body usecase.CreateNoteRequest true "Note"
// @Success 201 {object} usecase.NoteResponse
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/folders/{folderId}/notes [post]
func (h *AssetHandler) CreateNote(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	folderID, err := parseUintParam(c, "folderId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req usecase.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, usecaseErr := h.assetUsecase.CreateNote(userID, role, token, folderID, &req)
	if usecaseErr != nil {
		h.handleUsecaseError(c, usecaseErr)
		return
	}
	c.JSON(http.StatusCreated, note)
}

// ListNotesByFolder godoc
// @Summary List notes in folder
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param folderId path int true "Folder ID"
// @Success 200 {object} notesListBody
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/folders/{folderId}/notes [get]
func (h *AssetHandler) ListNotesByFolder(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	folderID, err := parseUintParam(c, "folderId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notes, usecaseErr := h.assetUsecase.ListNotesByFolder(userID, role, token, folderID)
	if usecaseErr != nil {
		h.handleUsecaseError(c, usecaseErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": notes})
}

// GetNote godoc
// @Summary Get note
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param noteId path int true "Note ID"
// @Success 200 {object} usecase.NoteResponse
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/notes/{noteId} [get]
func (h *AssetHandler) GetNote(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	noteID, err := parseUintParam(c, "noteId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, usecaseErr := h.assetUsecase.GetNote(userID, role, token, noteID)
	if usecaseErr != nil {
		h.handleUsecaseError(c, usecaseErr)
		return
	}
	c.JSON(http.StatusOK, note)
}

// UpdateNote godoc
// @Summary Update note
// @Tags notes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param noteId path int true "Note ID"
// @Param body body usecase.UpdateNoteRequest true "Note fields"
// @Success 200 {object} usecase.NoteResponse
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/notes/{noteId} [patch]
func (h *AssetHandler) UpdateNote(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	noteID, err := parseUintParam(c, "noteId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req usecase.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, usecaseErr := h.assetUsecase.UpdateNote(userID, role, token, noteID, &req)
	if usecaseErr != nil {
		h.handleUsecaseError(c, usecaseErr)
		return
	}
	c.JSON(http.StatusOK, note)
}

// DeleteNote godoc
// @Summary Delete note
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param noteId path int true "Note ID"
// @Success 200 {object} msgBody
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/notes/{noteId} [delete]
func (h *AssetHandler) DeleteNote(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	noteID, err := parseUintParam(c, "noteId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if usecaseErr := h.assetUsecase.DeleteNote(userID, role, token, noteID); usecaseErr != nil {
		h.handleUsecaseError(c, usecaseErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}

// ShareAsset godoc
// @Summary Share folder or note
// @Tags shares
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body usecase.ShareAssetRequest true "Share payload"
// @Success 200 {object} usecase.ShareResponse
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/shares [post]
func (h *AssetHandler) ShareAsset(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	var req usecase.ShareAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	share, err := h.assetUsecase.ShareAsset(userID, role, token, &req)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, share)
}

// RevokeShare godoc
// @Summary Revoke share
// @Tags shares
// @Security BearerAuth
// @Produce json
// @Param shareId path int true "Share ID"
// @Success 200 {object} msgBody
// @Failure 400 {object} errBody
// @Failure 403 {object} errBody
// @Failure 404 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/shares/{shareId} [delete]
func (h *AssetHandler) RevokeShare(c *gin.Context) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	shareID, err := parseUintParam(c, "shareId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if usecaseErr := h.assetUsecase.RevokeShare(userID, role, token, shareID); usecaseErr != nil {
		h.handleUsecaseError(c, usecaseErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Share revoked"})
}

// ListReceivedShares godoc
// @Summary Shares received by me
// @Tags shares
// @Security BearerAuth
// @Produce json
// @Success 200 {object} sharesListBody
// @Failure 401 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/shares/received [get]
func (h *AssetHandler) ListReceivedShares(c *gin.Context) {
	userID, _, _, ok := getActorContext(c)
	if !ok {
		return
	}

	shares, err := h.assetUsecase.ListReceivedShares(userID)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shares})
}

// ListGrantedShares godoc
// @Summary Shares granted by me
// @Tags shares
// @Security BearerAuth
// @Produce json
// @Success 200 {object} sharesListBody
// @Failure 401 {object} errBody
// @Failure 500 {object} errBody
// @Router /api/v1/shares/granted [get]
func (h *AssetHandler) ListGrantedShares(c *gin.Context) {
	userID, _, _, ok := getActorContext(c)
	if !ok {
		return
	}

	shares, err := h.assetUsecase.ListGrantedShares(userID)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shares})
}

func (h *AssetHandler) handleUsecaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrFolderNameRequired), errors.Is(err, usecase.ErrNoteTitleRequired),
		errors.Is(err, usecase.ErrInvalidAssetType), errors.Is(err, usecase.ErrInvalidAccessLevel),
		errors.Is(err, usecase.ErrCannotShareWithSelf):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrFolderNotFound), errors.Is(err, usecase.ErrNoteNotFound), errors.Is(err, usecase.ErrShareNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrShareTargetNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func parseUintParam(c *gin.Context, key string) (uint, error) {
	raw := c.Param(key)
	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid path parameter")
	}
	return uint(parsed), nil
}

func getActorContext(c *gin.Context) (uint, string, string, bool) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
		return 0, "", "", false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return 0, "", "", false
	}

	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role context"})
		return 0, "", "", false
	}

	role, ok := roleValue.(string)
	if !ok || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role context"})
		return 0, "", "", false
	}

	tokenValue, exists := c.Get("token")
	token := ""
	if exists {
		token, _ = tokenValue.(string)
	}

	return userID, role, token, true
}

type foldersListBody struct {
	Data []usecase.FolderResponse `json:"data"`
}

type notesListBody struct {
	Data []usecase.NoteResponse `json:"data"`
}

type sharesListBody struct {
	Data []usecase.ShareResponse `json:"data"`
}

type msgBody struct {
	Message string `json:"message"`
}

type errBody struct {
	Error string `json:"error"`
}
