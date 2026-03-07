package handlers

import (
	"church-system/internal/models"
	"church-system/internal/repository"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	admins        *repository.AdminRepository
	announcements *repository.AnnouncementRepository
	events        *repository.EventRepository
	gallery       *repository.GalleryRepository
	messages      *repository.MessageRepository
	uploadDir     string
}

func NewAdminHandler(
	adm *repository.AdminRepository,
	ann *repository.AnnouncementRepository,
	ev *repository.EventRepository,
	gal *repository.GalleryRepository,
	msg *repository.MessageRepository,
	uploadDir string,
) *AdminHandler {
	return &AdminHandler{
		admins:        adm,
		announcements: ann,
		events:        ev,
		gallery:       gal,
		messages:      msg,
		uploadDir:     uploadDir,
	}
}

// ─── Auth ─────────────────────────────────────────────────────────────────────

func (h *AdminHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_login.html", gin.H{"title": "Admin Login"})
}

func (h *AdminHandler) Login(c *gin.Context) {
	var form models.LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusOK, "admin_login.html", gin.H{
			"title": "Admin Login",
			"error": "Please enter valid email and password.",
		})
		return
	}

	admin, err := h.admins.FindByEmail(form.Email)
	if err != nil || admin == nil {
		c.HTML(http.StatusOK, "admin_login.html", gin.H{
			"title": "Admin Login",
			"error": "Invalid email or password.",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(form.Password)); err != nil {
		c.HTML(http.StatusOK, "admin_login.html", gin.H{
			"title": "Admin Login",
			"error": "Invalid email or password.",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("admin_id", admin.ID)
	session.Set("admin_name", admin.Name)
	session.Set("admin_email", admin.Email)
	session.Save()

	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func (h *AdminHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/admin/login")
}

// ─── Dashboard ────────────────────────────────────────────────────────────────

func (h *AdminHandler) Dashboard(c *gin.Context) {
	annCount, _ := h.announcements.Count()
	evCount, _ := h.events.Count()
	galCount, _ := h.gallery.Count()
	msgCount, _ := h.messages.CountUnread()

	recentAnn, _ := h.announcements.FindLatest(5)
	upcomingEv, _ := h.events.FindUpcoming(5)
	recentMsg, _ := h.messages.FindAll()
	if len(recentMsg) > 5 {
		recentMsg = recentMsg[:5]
	}

	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"title":         "Dashboard",
		"adminName":     session.Get("admin_name"),
		"annCount":      annCount,
		"evCount":       evCount,
		"galCount":      galCount,
		"msgCount":      msgCount,
		"recentAnn":     recentAnn,
		"upcomingEv":    upcomingEv,
		"recentMsg":     recentMsg,
		"activePage":    "dashboard",
	})
}

// ─── Announcements ────────────────────────────────────────────────────────────

func (h *AdminHandler) AnnouncementsList(c *gin.Context) {
	list, _ := h.announcements.FindAll(false)
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_announcements.html", gin.H{
		"title":         "Announcements",
		"announcements": list,
		"adminName":     session.Get("admin_name"),
		"activePage":    "announcements",
	})
}

func (h *AdminHandler) AnnouncementCreate(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_announcement_form.html", gin.H{
		"title":     "New Announcement",
		"adminName": session.Get("admin_name"),
		"activePage": "announcements",
	})
}

func (h *AdminHandler) AnnouncementStore(c *gin.Context) {
	var form models.AnnouncementForm
	if err := c.ShouldBind(&form); err != nil {
		session := sessions.Default(c)
		c.HTML(http.StatusOK, "admin_announcement_form.html", gin.H{
			"title":     "New Announcement",
			"error":     "Please fill all required fields.",
			"adminName": session.Get("admin_name"),
			"activePage": "announcements",
		})
		return
	}

	adminID := c.MustGet("admin_id").(int)
	category := form.Category
	if category == "" {
		category = "General"
	}

	if err := h.announcements.Create(form.Title, form.Content, category, form.IsPinned, adminID); err != nil {
		session := sessions.Default(c)
		c.HTML(http.StatusOK, "admin_announcement_form.html", gin.H{
			"title":     "New Announcement",
			"error":     "Failed to create announcement.",
			"adminName": session.Get("admin_name"),
			"activePage": "announcements",
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/announcements")
}

func (h *AdminHandler) AnnouncementEdit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ann, _ := h.announcements.FindByID(id)
	if ann == nil {
		c.Redirect(http.StatusFound, "/admin/announcements")
		return
	}
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_announcement_form.html", gin.H{
		"title":        "Edit Announcement",
		"announcement": ann,
		"adminName":    session.Get("admin_name"),
		"activePage":   "announcements",
	})
}

func (h *AdminHandler) AnnouncementUpdate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var form models.AnnouncementForm
	c.ShouldBind(&form)

	category := form.Category
	if category == "" {
		category = "General"
	}

	h.announcements.Update(id, form.Title, form.Content, category, form.IsPinned, form.IsActive)
	c.Redirect(http.StatusFound, "/admin/announcements")
}

func (h *AdminHandler) AnnouncementDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.announcements.Delete(id)
	c.Redirect(http.StatusFound, "/admin/announcements")
}

// ─── Events ───────────────────────────────────────────────────────────────────

func (h *AdminHandler) EventsList(c *gin.Context) {
	list, _ := h.events.FindAll(false)
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_events.html", gin.H{
		"title":      "Events",
		"events":     list,
		"adminName":  session.Get("admin_name"),
		"activePage": "events",
	})
}

func (h *AdminHandler) EventCreate(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_event_form.html", gin.H{
		"title":      "New Event",
		"adminName":  session.Get("admin_name"),
		"activePage": "events",
	})
}

func (h *AdminHandler) EventStore(c *gin.Context) {
	var form models.EventForm
	if err := c.ShouldBind(&form); err != nil {
		session := sessions.Default(c)
		c.HTML(http.StatusOK, "admin_event_form.html", gin.H{
			"title":      "New Event",
			"error":      "Please fill all required fields.",
			"adminName":  session.Get("admin_name"),
			"activePage": "events",
		})
		return
	}

	eventDate, err := time.Parse("2006-01-02", form.EventDate)
	if err != nil {
		session := sessions.Default(c)
		c.HTML(http.StatusOK, "admin_event_form.html", gin.H{
			"title":      "New Event",
			"error":      "Invalid date format.",
			"adminName":  session.Get("admin_name"),
			"activePage": "events",
		})
		return
	}

	adminID := c.MustGet("admin_id").(int)
	category := form.Category
	if category == "" {
		category = "General"
	}

	h.events.Create(form.Title, form.Description, form.Location, category,
		eventDate, form.StartTime, form.EndTime, form.IsRecurring, adminID)

	c.Redirect(http.StatusFound, "/admin/events")
}

func (h *AdminHandler) EventEdit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ev, _ := h.events.FindByID(id)
	if ev == nil {
		c.Redirect(http.StatusFound, "/admin/events")
		return
	}
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_event_form.html", gin.H{
		"title":      "Edit Event",
		"event":      ev,
		"adminName":  session.Get("admin_name"),
		"activePage": "events",
	})
}

func (h *AdminHandler) EventUpdate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var form models.EventForm
	c.ShouldBind(&form)

	eventDate, _ := time.Parse("2006-01-02", form.EventDate)
	category := form.Category
	if category == "" {
		category = "General"
	}

	h.events.Update(id, form.Title, form.Description, form.Location, category,
		eventDate, form.StartTime, form.EndTime, form.IsRecurring, form.IsActive)

	c.Redirect(http.StatusFound, "/admin/events")
}

func (h *AdminHandler) EventDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.events.Delete(id)
	c.Redirect(http.StatusFound, "/admin/events")
}

// ─── Gallery ──────────────────────────────────────────────────────────────────

func (h *AdminHandler) GalleryList(c *gin.Context) {
	list, _ := h.gallery.FindAll(false)
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_gallery.html", gin.H{
		"title":      "Gallery",
		"photos":     list,
		"adminName":  session.Get("admin_name"),
		"activePage": "gallery",
	})
}

func (h *AdminHandler) GalleryCreate(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_gallery_form.html", gin.H{
		"title":      "Upload Photo",
		"adminName":  session.Get("admin_name"),
		"activePage": "gallery",
	})
}

func (h *AdminHandler) GalleryStore(c *gin.Context) {
	session := sessions.Default(c)

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.HTML(http.StatusOK, "admin_gallery_form.html", gin.H{
			"title":      "Upload Photo",
			"error":      "Please select a photo to upload.",
			"adminName":  session.Get("admin_name"),
			"activePage": "gallery",
		})
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		c.HTML(http.StatusOK, "admin_gallery_form.html", gin.H{
			"title":      "Upload Photo",
			"error":      "Only image files (JPG, PNG, GIF, WebP) are allowed.",
			"adminName":  session.Get("admin_name"),
			"activePage": "gallery",
		})
		return
	}

	// Create unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	filename = strings.ReplaceAll(filename, " ", "_")

	// Ensure upload dir exists
	os.MkdirAll(h.uploadDir, 0755)
	savePath := filepath.Join(h.uploadDir, filename)

	out, err := os.Create(savePath)
	if err != nil {
		c.HTML(http.StatusOK, "admin_gallery_form.html", gin.H{
			"title":      "Upload Photo",
			"error":      "Failed to save photo.",
			"adminName":  session.Get("admin_name"),
			"activePage": "gallery",
		})
		return
	}
	defer out.Close()
	io.Copy(out, file)

	title := c.PostForm("title")
	description := c.PostForm("description")
	category := c.PostForm("category")
	if category == "" {
		category = "General"
	}

	adminID := c.MustGet("admin_id").(int)
	webPath := "/static/images/uploads/" + filename

	h.gallery.Create(title, description, filename, webPath, category, adminID)
	c.Redirect(http.StatusFound, "/admin/gallery")
}

func (h *AdminHandler) GalleryDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	photo, _ := h.gallery.FindByID(id)
	if photo != nil {
		// Delete physical file
		filePath := filepath.Join(h.uploadDir, photo.Filename)
		os.Remove(filePath)
		h.gallery.Delete(id)
	}
	c.Redirect(http.StatusFound, "/admin/gallery")
}

// ─── Messages ─────────────────────────────────────────────────────────────────

func (h *AdminHandler) MessagesList(c *gin.Context) {
	list, _ := h.messages.FindAll()
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_messages.html", gin.H{
		"title":      "Messages",
		"messages":   list,
		"adminName":  session.Get("admin_name"),
		"activePage": "messages",
	})
}

func (h *AdminHandler) MessageView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	msg, _ := h.messages.FindByID(id)
	if msg == nil {
		c.Redirect(http.StatusFound, "/admin/messages")
		return
	}
	h.messages.MarkRead(id)
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_message_view.html", gin.H{
		"title":      "Message",
		"message":    msg,
		"adminName":  session.Get("admin_name"),
		"activePage": "messages",
	})
}

func (h *AdminHandler) MessageDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.messages.Delete(id)
	c.Redirect(http.StatusFound, "/admin/messages")
}

// ─── Profile ──────────────────────────────────────────────────────────────────

func (h *AdminHandler) ProfilePage(c *gin.Context) {
	adminID := c.MustGet("admin_id").(int)
	admin, _ := h.admins.FindByID(adminID)
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "admin_profile.html", gin.H{
		"title":      "Profile",
		"admin":      admin,
		"adminName":  session.Get("admin_name"),
		"activePage": "profile",
	})
}

func (h *AdminHandler) ProfileUpdate(c *gin.Context) {
	adminID := c.MustGet("admin_id").(int)
	var form models.AdminProfileForm
	c.ShouldBind(&form)

	var hashedPw string
	if form.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err == nil {
			hashedPw = string(hashed)
		}
	}

	h.admins.Update(adminID, form.Name, form.Email, hashedPw)

	session := sessions.Default(c)
	session.Set("admin_name", form.Name)
	session.Set("admin_email", form.Email)
	session.Save()

	admin, _ := h.admins.FindByID(adminID)
	c.HTML(http.StatusOK, "admin_profile.html", gin.H{
		"title":      "Profile",
		"admin":      admin,
		"adminName":  session.Get("admin_name"),
		"success":    "Profile updated successfully.",
		"activePage": "profile",
	})
}
