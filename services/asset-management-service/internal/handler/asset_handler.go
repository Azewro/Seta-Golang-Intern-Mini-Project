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
