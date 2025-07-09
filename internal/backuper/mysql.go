package backuper

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type MySQLBackuper struct {
	Source DBSource
}

func (m *MySQLBackuper) Backup() (string, error) {

	src := m.Source
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.sql", src.GetTitle(), timestamp)
	filePath := filepath.Join("./backups", filename)
	gzPath := filePath + ".gz"

	// Ensure backups dir exists
	err := os.MkdirAll("./backups", 0755)
	if err != nil {
		return "", err
	}

	// Prepare dump command
	var cmd *exec.Cmd
	cmd = exec.Command("mysqldump",
		"-h", src.GetHost(),
		"-P", fmt.Sprint(src.GetPort()),
		"-u", src.GetUsername(),
		fmt.Sprintf("-p%s", src.GetPassword()),
		src.GetDBName(),
	)

	// Get the output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	// Open gzip file
	outFile, err := os.Create(gzPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
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
	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start dump: %w", err)
	}

	// Copy dump into gzip file
	if _, err = io.Copy(gzWriter, stdout); err != nil {
		return "", fmt.Errorf("failed to copy dump to gzip: %w", err)
	}

	if err = cmd.Wait(); err != nil {
		return "", fmt.Errorf("dump command failed: %w", err)
	}

	return gzPath, nil
}
