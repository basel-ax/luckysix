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
		// Lock file exists, check if it's stale
		if time.Since(info.ModTime()) > LockFileTimeout {
			// Lock is stale, remove it
			log.Printf("Removing stale lock for %s", c.commandName)
			if err := os.Remove(c.lockFile); err != nil {
				return false, fmt.Errorf("failed to remove stale lock: %w", err)
			}
		} else {
			// Lock is still active
			log.Printf("%s is already running (lock file exists)", c.commandName)
			return false, nil
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
	defer lock.Release()

	// Run the command
	if err := runFunc(); err != nil {
		return true, err
	}

	return true, nil
}
