package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	// LockFileDir is the directory where lock files are stored
	LockFileDir = "/tmp/luckysix"
	// LockFileTimeout is the maximum time a lock can exist before it's considered stale
	LockFileTimeout = 24 * time.Hour
)

// CronLock represents a lock for a specific command
type CronLock struct {
	commandName string
	lockFile    string
}

// getPIDFromLockFile extracts the PID from the lock file
func (c *CronLock) getPIDFromLockFile() (int, error) {
	data, err := os.ReadFile(c.lockFile)
	if err != nil {
		return 0, fmt.Errorf("failed to read lock file: %w", err)
	}

	// Parse PID from first line: "PID: 12345"
	var pid int
	_, err = fmt.Sscanf(string(data), "PID: %d", &pid)
	if err != nil {
		return 0, fmt.Errorf("failed to parse PID from lock file: %w", err)
	}

	return pid, nil
}

// isProcessRunning checks if a process with the given PID is currently running
func (c *CronLock) isProcessRunning(pid int) bool {
	// On Unix-like systems, we can check if a process exists by looking at /proc/<pid>
	// This approach works without needing to send signals
	if pid <= 0 {
		return false
	}

	// Check if /proc/<pid> exists
	_, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	return err == nil
}

// NewCronLock creates a new CronLock for the given command
func NewCronLock(commandName string) *CronLock {
	// Ensure lock directory exists
	if err := os.MkdirAll(LockFileDir, 0755); err != nil {
		log.Printf("Warning: Could not create lock directory: %v", err)
	}

	return &CronLock{
		commandName: commandName,
		lockFile:    filepath.Join(LockFileDir, commandName+".lock"),
	}
}

// Acquire tries to acquire a lock for the command
// Returns true if lock was acquired, false if already running
func (c *CronLock) Acquire() (bool, error) {
	// Check if lock file exists
	info, err := os.Stat(c.lockFile)
	if err == nil {
		// Lock file exists, check if it's stale or if process is still running
		if time.Since(info.ModTime()) > LockFileTimeout {
			// Lock is stale, remove it
			log.Printf("Removing stale lock for %s", c.commandName)
			if err := os.Remove(c.lockFile); err != nil {
				return false, fmt.Errorf("failed to remove stale lock: %w", err)
			}
		} else {
			// Check if the process is still running
			if pid, err := c.getPIDFromLockFile(); err == nil && pid > 0 {
				if c.isProcessRunning(pid) {
					// Lock is still active and process is running
					log.Printf("%s is already running (PID: %d)", c.commandName, pid)
					return false, nil
				}
				// Process is not running, remove stale lock
				log.Printf("Removing stale lock for %s (PID %d not running)", c.commandName, pid)
				if err := os.Remove(c.lockFile); err != nil {
					return false, fmt.Errorf("failed to remove stale lock: %w", err)
				}
			} else {
				// Could not read PID from lock file, treat as stale
				log.Printf("Removing stale lock for %s (could not read PID)", c.commandName)
				if err := os.Remove(c.lockFile); err != nil {
					return false, fmt.Errorf("failed to remove stale lock: %w", err)
				}
			}
		}
	} else if !os.IsNotExist(err) {
		// Error checking lock file
		return false, fmt.Errorf("failed to check lock file: %w", err)
	}

	// Create lock file
	f, err := os.Create(c.lockFile)
	if err != nil {
		return false, fmt.Errorf("failed to create lock file: %w", err)
	}
	defer f.Close()

	// Write PID and timestamp to lock file
	_, err = fmt.Fprintf(f, "PID: %d\nStarted: %s\n", os.Getpid(), time.Now().Format(time.RFC3339))
	if err != nil {
		return false, fmt.Errorf("failed to write lock file: %w", err)
	}

	log.Printf("Acquired lock for %s", c.commandName)
	return true, nil
}

// Release removes the lock file
func (c *CronLock) Release() error {
	if err := os.Remove(c.lockFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove lock file: %w", err)
	}
	log.Printf("Released lock for %s", c.commandName)
	return nil
}

// CheckAndRunCommand checks if the command is already running and runs it if not
// Returns true if command was run, false if already running
func CheckAndRunCommand(commandName string, runFunc func() error) (bool, error) {
	lock := NewCronLock(commandName)

	acquired, err := lock.Acquire()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !acquired {
		return false, nil
	}

	// Ensure lock is released when done
	defer func() {
		if err := lock.Release(); err != nil {
			log.Printf("Warning: failed to release lock for %s: %v", commandName, err)
		}
	}()

	// Run the command
	if err := runFunc(); err != nil {
		return true, err
	}

	return true, nil
}
