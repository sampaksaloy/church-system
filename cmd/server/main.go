package main

import (
	"church-system/internal/config"
	"church-system/internal/database"
	"church-system/internal/handlers"
	"church-system/internal/middleware"
	"church-system/internal/repository"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	loadDotEnv(".env")
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	adminRepo := repository.NewAdminRepository(db)
	count, _ := adminRepo.Count()
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(cfg.AdminDefaultPassword), bcrypt.DefaultCost)
		adminRepo.Create("Church Administrator", cfg.AdminDefaultEmail, string(hashed))
		log.Printf("✅ Default admin created: %s / %s", cfg.AdminDefaultEmail, cfg.AdminDefaultPassword)
	}

	os.MkdirAll(cfg.UploadDir, 0755)

	r := gin.Default()

	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string { return t.Format("January 2, 2006") },
		"formatTime": func(s string) string {
			if s == "" { return "" }
			t, err := time.Parse("15:04:05", s)
			if err != nil { t, err = time.Parse("15:04", s); if err != nil { return s } }
			return t.Format("3:04 PM")
		},
		"formatDateTime": func(t time.Time) string { return t.Format("Jan 2, 2006 3:04 PM") },
		"shortDate":      func(t time.Time) string { return t.Format("Jan 2") },
		"isoDate":        func(t time.Time) string { return t.Format("2006-01-02") },
		"truncate": func(s string, n int) string {
			if len(s) <= n { return s }
			return s[:n] + "..."
		},
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"mod": func(a, b int) int { return a % b },
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := range s { s[i] = i + 1 }
			return s
		},
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"nl2br": func(s string) template.HTML {
			return template.HTML(strings.ReplaceAll(template.HTMLEscapeString(s), "\n", "<br>"))
		},
		"slice": func(items ...string) []string { return items },
		"now":   func() time.Time { return time.Now() },
	}

	// Load templates from all subdirectories
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/partials/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/public/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/admin/*.html"))
	r.SetHTMLTemplate(tmpl)

	r.Static("/static", "./static")

	store := cookie.NewStore([]byte(cfg.SessionSecret))
	r.Use(sessions.Sessions("church_session", store))

	annRepo := repository.NewAnnouncementRepository(db)
	evRepo  := repository.NewEventRepository(db)
	galRepo := repository.NewGalleryRepository(db)
	msgRepo := repository.NewMessageRepository(db)

	pub := handlers.NewPublicHandler(annRepo, evRepo, galRepo, msgRepo)
	adm := handlers.NewAdminHandler(adminRepo, annRepo, evRepo, galRepo, msgRepo, cfg.UploadDir)

	r.GET("/", pub.Home)
	r.GET("/announcements", pub.Announcements)
	r.GET("/announcements/:id", pub.AnnouncementDetail)
	r.GET("/events", pub.Events)
	r.GET("/gallery", pub.Gallery)
	r.GET("/contact", pub.Contact)
	r.POST("/contact", pub.ContactSubmit)

	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/login", middleware.GuestOnly(), adm.LoginPage)
		adminGroup.POST("/login", middleware.GuestOnly(), adm.Login)

		protected := adminGroup.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/logout", adm.Logout)
			protected.GET("/dashboard", adm.Dashboard)

			protected.GET("/announcements", adm.AnnouncementsList)
			protected.GET("/announcements/create", adm.AnnouncementCreate)
			protected.POST("/announcements/create", adm.AnnouncementStore)
			protected.GET("/announcements/:id/edit", adm.AnnouncementEdit)
			protected.POST("/announcements/:id/edit", adm.AnnouncementUpdate)
			protected.POST("/announcements/:id/delete", adm.AnnouncementDelete)

			protected.GET("/events", adm.EventsList)
			protected.GET("/events/create", adm.EventCreate)
			protected.POST("/events/create", adm.EventStore)
			protected.GET("/events/:id/edit", adm.EventEdit)
			protected.POST("/events/:id/edit", adm.EventUpdate)
			protected.POST("/events/:id/delete", adm.EventDelete)

			protected.GET("/gallery", adm.GalleryList)
			protected.GET("/gallery/upload", adm.GalleryCreate)
			protected.POST("/gallery/upload", adm.GalleryStore)
			protected.POST("/gallery/:id/delete", adm.GalleryDelete)

			protected.GET("/messages", adm.MessagesList)
			protected.GET("/messages/:id", adm.MessageView)
			protected.POST("/messages/:id/delete", adm.MessageDelete)

			protected.GET("/profile", adm.ProfilePage)
			protected.POST("/profile", adm.ProfileUpdate)
		}
	}

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("🚀 Church System running at http://%s", addr)
	log.Printf("📋 Admin panel: http://%s/admin/login", addr)
	r.Run(addr)
}

func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil { return }
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") { continue }
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" { os.Setenv(key, val) }
		}
	}
}
