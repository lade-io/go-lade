package lade

import (
	"time"

	"github.com/saulortega/pgeo.latlng"
)

type Addon struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Owner        *User     `json:"owner"`
	Service      *Service  `json:"service"`
	PlanID       string    `json:"plan_id"`
	Region       *Region   `json:"region"`
	Release      string    `json:"release"`
	Public       bool      `json:"public"`
	Status       string    `json:"status"`
	Hostname     string    `json:"hostname"`
	Port         int       `json:"port"`
	Database     string    `json:"database"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	BackupLimit  int       `json:"backup_limit"`
	BackupWindow string    `json:"backup_window"`
	CreatedAt    time.Time `json:"created_at"`
}

type App struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Owner     *User     `json:"owner"`
	Region    *Region   `json:"region"`
	Status    string    `json:"status"`
	Hostname  string    `json:"hostname"`
	CreatedAt time.Time `json:"created_at"`
}

type Attachment struct {
	ID        int       `json:"id"`
	AppID     int       `json:"app_id"`
	AddonID   int       `json:"addon_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Backup struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	AddonID   int       `json:"addon_id"`
	Resource  string    `json:"resource"`
	SizeGB    int       `json:"size_gb"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Container struct {
	ID        string    `json:"id"`
	AppID     int       `json:"app_id"`
	PlanID    string    `json:"plan_id"`
	Process   *Process  `json:"process"`
	Number    int       `json:"number"`
	CreatedAt time.Time `json:"created_at"`
}

type Domain struct {
	ID        int       `json:"id"`
	AppID     int       `json:"app_id"`
	Hostname  string    `json:"hostname"`
	CreatedAt time.Time `json:"created_at"`
}

type Env struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LogEntry struct {
	Line   string `json:"line"`
	Name   string `json:"name"`
	Source string `json:"source"`
}

type Plan struct {
	ID           string  `json:"id"`
	Cpu          int     `json:"cpu"`
	Ram          int     `json:"ram"`
	Disk         int     `json:"disk"`
	PriceHourly  float64 `json:"price_hourly"`
	PriceMonthly float64 `json:"price_monthly"`
}

type Process struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	AppID     int       `json:"app_id"`
	ReleaseID int       `json:"release_id"`
	PlanID    string    `json:"plan_id"`
	Command   string    `json:"command"`
	Count     int       `json:"count"`
	Number    int       `json:"number"`
	Replicas  int       `json:"replicas"`
	CreatedAt time.Time `json:"created_at"`
}

type Quota struct {
	ID     int    `json:"id"`
	PlanID string `json:"plan_id"`
	UserID int    `json:"user_id"`
	Quota  int    `json:"quota"`
}

type Region struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Country  string     `json:"country"`
	Location pgeo.Point `json:"location"`
}

type Release struct {
	ID        int       `json:"id"`
	Version   int       `json:"version"`
	AppID     int       `json:"app_id"`
	RepoID    int       `json:"repo_id"`
	Active    bool      `json:"active"`
	Branch    string    `json:"branch"`
	Commit    string    `json:"commit"`
	CreatedAt time.Time `json:"created_at"`
}

type Repo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	OwnerID   int       `json:"owner_id"`
	GitURL    string    `json:"git_url"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

type Service struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Repo      *Repo     `json:"repo"`
	Connector string    `json:"connector"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	RegionID  string    `json:"region_id"`
	CreatedAt time.Time `json:"created_at"`
}
