// storage/session.go
package storage

import (
	"fmt"
)

type Session struct {
	ID      int `gorm:"primary_key"`
	Version int
	Phone   string `gorm:"uniqueIndex"`
	Data    []byte
}

// For future migrations
// type Session1 struct {
//     Version   int `gorm:"primary_key"`
//     DC        int
//     Addr      string
//     AuthKey   []byte
//     AuthKeyID []byte
// }

const LatestVersion = 1

// UpdateSession saves or updates the session in storage
func (ps *PeerStorage) UpdateSession(session *Session) error {
	if session.Phone == "" {
		return fmt.Errorf("phone number is required")
	}

	tx := ps.SqlSession.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Try to find existing session first
	var existing Session

	if err := tx.Where("phone = ?", session.Phone).First(&existing).Error; err != nil && !IsNotFoundError(err) {
		tx.Rollback()
		return fmt.Errorf("check existing session: %w", err)
	}

	// If session exists, update it, otherwise create new
	if err := tx.Save(session).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("save session: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// GetSession returns the session saved in storage.
func (ps *PeerStorage) GetSession(phone string) (*Session, error) {
	if phone == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	session := Session{
		Version: LatestVersion,
		Phone:   phone,
	}

	if err := ps.SqlSession.Where("phone = ?", phone).First(&session).Error; err != nil {
		if IsNotFoundError(err) {
			return &session, nil
		}
		return nil, fmt.Errorf("get session: %w", err)
	}

	return &session, nil
}

// GetAllSessions returns all sessions from storage
func (ps *PeerStorage) GetAllSessions() ([]*Session, error) {
	var sessions []*Session

	if err := ps.SqlSession.Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("get all sessions: %w", err)
	}

	return sessions, nil
}

// DeleteSession removes a session from storage
func (ps *PeerStorage) DeleteSession(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number is required")
	}

	tx := ps.SqlSession.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if err := tx.Where("phone = ?", phone).Delete(&Session{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("delete session: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Helper function to check if error is a "not found" error
func IsNotFoundError(err error) bool {
	return err.Error() == "record not found"
}
