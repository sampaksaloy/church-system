package handlers

import (
	"church-system/internal/models"
	"church-system/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	announcements *repository.AnnouncementRepository
	events        *repository.EventRepository
	gallery       *repository.GalleryRepository
	messages      *repository.MessageRepository
}

func NewPublicHandler(
	ann *repository.AnnouncementRepository,
	ev *repository.EventRepository,
	gal *repository.GalleryRepository,
	msg *repository.MessageRepository,
) *PublicHandler {
	return &PublicHandler{
		announcements: ann,
		events:        ev,
		gallery:       gal,
		messages:      msg,
	}
}

func (h *PublicHandler) Home(c *gin.Context) {
	announcements, _ := h.announcements.FindLatest(5)
	events, _ := h.events.FindUpcoming(6)
	photos, _ := h.gallery.FindLatest(6)

	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":         "Home",
		"announcements": announcements,
		"events":        events,
		"photos":        photos,
		"currentPage":   "home",
		"year":          time.Now().Year(),
	})
}

func (h *PublicHandler) Announcements(c *gin.Context) {
	category := c.Query("category")
	all, _ := h.announcements.FindAll(true)

	var filtered []*models.Announcement
	for _, a := range all {
		if category == "" || a.Category == category {
			filtered = append(filtered, a)
		}
	}

	c.HTML(http.StatusOK, "announcements.html", gin.H{
		"title":         "Announcements",
		"announcements": filtered,
		"category":      category,
		"currentPage":   "announcements",
		"year":          time.Now().Year(),
	})
}

func (h *PublicHandler) AnnouncementDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/announcements")
		return
	}
	a, _ := h.announcements.FindByID(id)
	if a == nil || !a.IsActive {
		c.Redirect(http.StatusFound, "/announcements")
		return
	}
	c.HTML(http.StatusOK, "announcement_detail.html", gin.H{
		"title":        a.Title,
		"announcement": a,
		"currentPage":  "announcements",
		"year":         time.Now().Year(),
	})
}

func (h *PublicHandler) Events(c *gin.Context) {
	// Month navigation
	now := time.Now()
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	year := now.Year()
	month := int(now.Month())

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}
	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = m
		}
	}

	events, _ := h.events.FindByMonth(year, month)
	upcoming, _ := h.events.FindUpcoming(10)

	// Build calendar grid
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	// Map events to days
	eventMap := make(map[int][]*models.Event)
	for _, e := range events {
		day := e.EventDate.Day()
		eventMap[day] = append(eventMap[day], e)
	}

	// Prev/next month
	prevMonth := firstDay.AddDate(0, -1, 0)
	nextMonth := firstDay.AddDate(0, 1, 0)

	c.HTML(http.StatusOK, "events.html", gin.H{
		"title":        "Event Calendar",
		"events":       events,
		"upcoming":     upcoming,
		"eventMap":     eventMap,
		"year":         year,
		"month":        month,
		"monthName":    firstDay.Format("January"),
		"firstWeekday": int(firstDay.Weekday()),
		"daysInMonth":  lastDay.Day(),
		"prevYear":     prevMonth.Year(),
		"prevMonth":    int(prevMonth.Month()),
		"nextYear":     nextMonth.Year(),
		"nextMonth":    int(nextMonth.Month()),
		"currentPage":  "events",
	})
}

func (h *PublicHandler) Gallery(c *gin.Context) {
	category := c.Query("category")
	all, _ := h.gallery.FindAll(true)

	var filtered []*models.GalleryPhoto
	for _, p := range all {
		if category == "" || p.Category == category {
			filtered = append(filtered, p)
		}
	}

	c.HTML(http.StatusOK, "gallery.html", gin.H{
		"title":       "Photo Gallery",
		"photos":      filtered,
		"category":    category,
		"currentPage": "gallery",
		"year":        time.Now().Year(),
	})
}

func (h *PublicHandler) Contact(c *gin.Context) {
	c.HTML(http.StatusOK, "contact.html", gin.H{
		"title":       "Contact Us",
		"currentPage": "contact",
		"year":        time.Now().Year(),
	})
}

func (h *PublicHandler) ContactSubmit(c *gin.Context) {
	var form models.ContactForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusOK, "contact.html", gin.H{
			"title":       "Contact Us",
			"error":       "Please fill all required fields correctly.",
			"currentPage": "contact",
		})
		return
	}

	if err := h.messages.Create(form.SenderName, form.Email, form.Phone, form.Subject, form.Message); err != nil {
		c.HTML(http.StatusOK, "contact.html", gin.H{
			"title":       "Contact Us",
			"error":       "Failed to send message. Please try again.",
			"currentPage": "contact",
		})
		return
	}

	c.HTML(http.StatusOK, "contact.html", gin.H{
		"title":       "Contact Us",
		"success":     "Your message has been sent! We will get back to you soon.",
		"currentPage": "contact",
		"year":        time.Now().Year(),
	})
}
