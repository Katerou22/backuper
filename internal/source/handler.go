package source

import (
	"compress/gzip"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
func (h *Handler) CreateFromDsn(dsn string) error {

	src, err := parseDSN(dsn)
	if err != nil {
		return err
	}

	err = h.repo.Create(src)
	if err != nil {
		return nil
	}

	return nil
}

func (h *Handler) Delete(src *Source) error {

	err := h.repo.FindOne(src)
	if err != nil {
		return err
	}

	return h.repo.Delete(src)

}

func (h *Handler) Backup(src *Source) error {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.sql", src.Title, timestamp)
	filePath := filepath.Join("./backups", filename)
	gzPath := filePath + ".gz"

	// Ensure backups dir exists
	err := os.MkdirAll("./backups", 0755)
	if err != nil {
		return err
	}

	// Prepare dump command
	var cmd *exec.Cmd
	switch src.DBType {
	case "mysql":
		cmd = exec.Command("mysqldump",
			"-h", src.Host,
			"-P", fmt.Sprint(src.Port),
			"-u", src.Username,
			fmt.Sprintf("-p%s", src.Password),
			src.DB,
		)
	case "postgres":
		_ = os.Setenv("PGPASSWORD", src.Password)
		cmd = exec.Command("pg_dump",
			"-h", src.Host,
			"-p", fmt.Sprint(src.Port),
			"-U", src.Username,
			"-d", src.DB,
		)
	default:
		return fmt.Errorf("unsupported db type: %s", src.DBType)
	}

	// Get the output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	// Open gzip file
	outFile, err := os.Create(gzPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {

		}
	}(outFile)

	gzWriter := gzip.NewWriter(outFile)
	defer func(gzWriter *gzip.Writer) {
		err := gzWriter.Close()
		if err != nil {

		}
	}(gzWriter)

	// Start command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start dump: %w", err)
	}

	// Copy dump into gzip file
	if _, err := io.Copy(gzWriter, stdout); err != nil {
		return fmt.Errorf("failed to copy dump to gzip: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("dump command failed: %w", err)
	}

	return nil
}

func parseDSN(dsn string) (*Source, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DSN: %w", err)
	}

	pw, _ := u.User.Password()
	pwDecoded, _ := url.QueryUnescape(pw)

	port, _ := strconv.Atoi(u.Port())

	dbType := u.Scheme // "postgres", "mysql", etc.

	return &Source{
		Link:     dsn,
		Host:     u.Hostname(),
		Port:     port,
		DB:       strings.TrimPrefix(u.Path, "/"),
		Username: u.User.Username(),
		Password: pwDecoded,
		DBType:   dbType,
	}, nil
}
