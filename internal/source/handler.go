package source

import (
	"backuper/internal/backuper"
	"fmt"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"
)

type Handler struct {
	repo *Repository
}

func NewHandler(db *gorm.DB) (*Handler, error) {
	repo := NewRepository(db)

	err := repo.Migrate()

	if err != nil {
		return nil, err
	}

	return &Handler{repo: repo}, nil
}
func (h *Handler) CreateFromDsn(title string, dsn string) error {

	src, err := parseDSN(dsn)
	if err != nil {
		return err
	}

	src.Title = title

	err = h.repo.Create(src)
	if err != nil {
		return nil
	}

	return nil
}

func (h *Handler) List() []*Source {

	var sources []*Source
	h.repo.FindAll(&sources)

	return sources
}

func (h *Handler) Delete(src *Source) error {

	err := h.repo.FindOne(src)
	if err != nil {
		return err
	}

	return h.repo.Delete(src)

}

func (h *Handler) Find(id uint) *Source {
	src, _ := h.repo.FindByID(id)

	return src
}

func (h *Handler) Backup(id uint) (string, error) {
	src, err := h.repo.FindByID(id)

	if err != nil {
		return "", err
	}
	//
	b, err := backuper.NewBackuper(src)

	if err != nil {
		return "", err
	}

	path, err := b.Backup()

	if err != nil {
		return "", err
	}

	fmt.Println(path)
	return path, nil

}

func (h *Handler) BackupAll() {
	list := h.List()

	if len(list) > 0 {
		for _, src := range list {
			h.Backup(src.ID)
		}
	}

}

func parseDSN(dsn string) (*Source, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DSN: %w", err)
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("missing database type (scheme) (mysql or postgres)")
	}

	if u.Hostname() == "" {
		return nil, fmt.Errorf("missing host")
	}

	portStr := u.Port()
	if portStr == "" {
		return nil, fmt.Errorf("missing port")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %s", portStr)
	}

	if u.User == nil || u.User.Username() == "" {
		return nil, fmt.Errorf("missing username")
	}

	pw, hasPassword := u.User.Password()
	if !hasPassword || pw == "" {
		return nil, fmt.Errorf("missing password")
	}
	pwDecoded, err := url.QueryUnescape(pw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode password: %w", err)
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		return nil, fmt.Errorf("missing database name")
	}

	return &Source{
		Link:     dsn,
		Host:     u.Hostname(),
		Port:     port,
		DB:       dbName,
		Username: u.User.Username(),
		Password: pwDecoded,
		DBType:   u.Scheme,
	}, nil
}
