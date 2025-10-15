package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/devhammed/grey-schedule/backend/internal/store"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	store  store.Store
}

type createReq struct {
	Title string `json:"title"`
	Start string `json:"start"`
	End   string `json:"end"`
}

func NewServer(s store.Store) *Server {
	srv := &Server{engine: gin.Default(), store: s}

	srv.cors()

	srv.routes()

	return srv
}

func (s *Server) cors() {
	s.engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

func (s *Server) routes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	{
		api := s.engine.Group("/api")

		api.GET("/appointments", s.listAppointments)

		api.POST("/appointments", s.createAppointment)

		api.DELETE("/appointments/:id", s.deleteAppointment)
	}
}

func (s *Server) listAppointments(c *gin.Context) {
	list := s.store.List()

	c.JSON(http.StatusOK, list)
}

func (s *Server) createAppointment(c *gin.Context) {
	var req createReq

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}

	start, err1 := time.Parse(time.RFC3339, req.Start)

	end, err2 := time.Parse(time.RFC3339, req.End)

	if err1 != nil || err2 != nil {
		writeError(c, http.StatusBadRequest, "start and end must be RFC3339 timestamps")
		return
	}

	appointment, err := s.store.Create(req.Title, start, end)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			writeError(c, http.StatusConflict, err.Error())
		default:
			writeError(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

func (s *Server) deleteAppointment(c *gin.Context) {
	id := c.Param("id")

	if err := s.store.Delete(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(c, http.StatusNotFound, err.Error())
			return
		}

		writeError(c, http.StatusInternalServerError, "internal error")

		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) Start(addr string) error {
	return s.engine.Run(addr)
}

func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}
