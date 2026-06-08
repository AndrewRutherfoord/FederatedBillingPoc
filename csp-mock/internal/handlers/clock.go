package handlers

import "github.com/gin-gonic/gin"

type advanceTimeRequest struct {
	Seconds int64 `json:"seconds"`
}

// GetCurrentTime godoc
//
//	@Summary	Get the current time from the mock clock
//	@Tags		clock
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Failure	401	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/clock/current [get]
func (s *Server) GetCurrentTime(c *gin.Context) {
	currentTime := s.clock.Now().Format("2006-01-02T15:04:05Z")
	c.JSON(200, gin.H{"current_time": currentTime})
}

// AdvanceTime godoc
//
//	@Summary	Advance the mock clock by a specified number of seconds. Default to 1 hour if not provided.
//	@Tags		clock
//	@Accept		json
//	@Produce	json
//	@Param		body	body advanceTimeRequest true	"Number of seconds to advance"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	401	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/clock/advance [post]
func (s *Server) AdvanceTime(c *gin.Context) {
	var req advanceTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Seconds <= 0 {
		req.Seconds = 3600 // Default to 1 hour
	}

	adv, ok := s.clock.(interface{ Advance(int64) })
	if !ok {
		c.JSON(501, gin.H{"error": "clock does not support time advancement"})
		return
	}
	adv.Advance(req.Seconds)

	// Return the new current time after advancement
	s.GetCurrentTime(c)
}
