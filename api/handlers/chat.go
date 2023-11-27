package handlers

import (
	"context"
	"fmt"
	"golang-websocker/api/http"
	"golang-websocker/pkg/models"
	"golang-websocker/pkg/wsclient"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Chat godoc
// @Summary Chat handler
// @Description you can chat with your friends
// @Tags chat
// @Accept  json
// @Produce  json
// @Param chat_id path string true "chat_id"
// @Param sender_name query string true "sender_name"
// @Param user_id path string true "user_id"
// @Param body body models.Message true "body"
// @Success 200 {object} http.Response{data=string} "ok"
// @Failure 400 {object} http.Response{data=string} "bad request"
// @Failure 500 {object} http.Response{data=string} "internal server error"
// @Router /ws/{chat_id}/{user_id} [get]
func (h *Handler) WsHandler(c *gin.Context) {

	chat_id := c.Param("chat_id")
	if _, err := uuid.Parse(chat_id); err != nil {
		h.handleResponse(c, http.BadRequest, fmt.Errorf("invalid chat_id: %s", err))
		return
	}

	userID := c.Param("user_id")
	senderName := c.Query("sender_name")

	data := models.CtxData{
		UserID:     userID,
		SenderName: senderName,
	}

	ctx := context.WithValue(c.Request.Context(), "data", data)

	wsclient.ServeWS(ctx, c.Writer, c.Request, chat_id, h.hub)
}
