package handler

import (
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createList(c *gin.Context) {
	h.Logger.Println("CREATE LIST")
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input model.TodoListDTO
	h.Logger.Println("DECODE MODEL TODOLIST")
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	h.Logger.Println("create list handler!")
	id, err := h.Service.TodoList.Create(c.Request.Context(), userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllListsResponse struct {
	Data []model.TodoListDTO `json:"data"`
}

func (h *Handler) getAllList(c *gin.Context) {
	h.Logger.Println("GET ALL LIST")

	h.Logger.Println("GET ID")
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.Service.TodoList.GetAll(c.Request.Context(), userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

func (h *Handler) getListById(c *gin.Context) {
	h.Logger.Println("GET LIST BY ID")
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	h.Logger.Println("GET ID FROM CONTEXT")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	list, err := h.Service.TodoList.GetById(c.Request.Context(), userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
	h.Logger.Println("UPDATE LIST")

	h.Logger.Println("GET USER ID")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Logger.Println("GET LIST ID")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	h.Logger.Println("DECODE MODEL")
	var input model.UpdateListInputDTO
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Service.TodoList.Update(c.Request.Context(), userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteList(c *gin.Context) {
	h.Logger.Println("DELETE LIST")
	h.Logger.Println("GET USER ID")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.Logger.Println("GET ID LIST")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.Service.TodoList.Delete(c.Request.Context(), userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
