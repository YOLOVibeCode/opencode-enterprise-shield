// Package audit provides audit logging for Enterprise Shield.
package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
	"github.com/google/uuid"
)

// Logger provides audit logging functionality.
type Logger struct {
	logPath           string
	signEntries       bool
	signer            *Signer
	lastEntryHash     string
	mu                sync.Mutex
	file              *os.File
	retentionDays     int
}

// NewLogger creates a new audit logger.
func NewLogger(logPath string, signEntries bool, retentionDays int) (*Logger, error) {
	// Expand home directory
	if logPath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		logPath = filepath.Join(home, logPath[1:])
	}

	// Create log directory
	if err := os.MkdirAll(logPath, 0750); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file for today
	logFile := filepath.Join(logPath, fmt.Sprintf("audit_%s.jsonl", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := &Logger{
		logPath:       logPath,
		signEntries:   signEntries,
		retentionDays: retentionDays,
		file:          file,
	}

	if signEntries {
		logger.signer = NewSigner()
	}

	return logger, nil
}

// Log writes an audit entry asynchronously.
func (l *Logger) Log(entry types.AuditEntry) {
	go func() {
		_ = l.logEntry(entry)
	}()
}

// LogSync writes an audit entry synchronously.
func (l *Logger) LogSync(entry types.AuditEntry) error {
	return l.logEntry(entry)
}

// logEntry is the internal logging function.
func (l *Logger) logEntry(entry types.AuditEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Set entry ID and timestamp if not set
	if entry.EntryID == "" {
		entry.EntryID = "audit_" + uuid.New().String()[:12]
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	// Compute request hash if not set
	if entry.RequestHash == "" {
		entry.RequestHash = computeHash(entry.EntryID + entry.UserID + entry.Timestamp.String())
	}

	// Set previous entry hash for chain integrity
	entry.PreviousEntryHash = l.lastEntryHash

	// Sign entry if signing is enabled
	if l.signEntries && l.signer != nil {
		signature, err := l.signer.Sign(entry)
		if err == nil {
			entry.Signature = signature
		}
	}

	// Serialize entry
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal audit entry: %w", err)
	}

	// Write to file
	if _, err := l.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write audit entry: %w", err)
	}

	// Update last entry hash
	l.lastEntryHash = computeHash(string(data))

	return nil
}

// CreateEntry creates a new audit entry from request/response data.
func (l *Logger) CreateEntry(
	userID string,
	sessionID string,
	department string,
	provider string,
	wasSanitized bool,
	violations []types.Violation,
	action types.Action,
	processingTimeMs int64,
) types.AuditEntry {
	return types.AuditEntry{
		EntryID:          "audit_" + uuid.New().String()[:12],
		Timestamp:        time.Now().UTC(),
		UserID:           userID,
		SessionID:        sessionID,
		Department:       department,
		Provider:         provider,
		WasSanitized:     wasSanitized,
		Violations:       violations,
		Action:           action,
		ProcessingTimeMs: processingTimeMs,
	}
}

// RotateFile rotates the log file (creates new file for new day).
func (l *Logger) RotateFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Close current file
	if l.file != nil {
		l.file.Close()
	}

	// Create new log file
	logFile := filepath.Join(l.logPath, fmt.Sprintf("audit_%s.jsonl", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	l.file = file
	return nil
}

// CleanupOldLogs removes log files older than retention period.
func (l *Logger) CleanupOldLogs() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -l.retentionDays)

	entries, err := os.ReadDir(l.logPath)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(l.logPath, entry.Name()))
		}
	}

	return nil
}

// Close closes the logger.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// computeHash computes SHA256 hash of a string.
func computeHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

