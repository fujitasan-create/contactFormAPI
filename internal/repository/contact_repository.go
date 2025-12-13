package repository

import (
	"time"

	"contactFormAPI/internal/db"
)

type Contact struct {
	ID        int64     `json:"id"`
	Contact   string    `json:"contact"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	IP        *string   `json:"ip,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty"`
}

type ContactRepository struct{}

func NewContactRepository() *ContactRepository {
	return &ContactRepository{}
}

func (r *ContactRepository) Create(contact, name, message string, ip, userAgent *string) (*Contact, error) {
	query := `
		INSERT INTO contacts (contact, name, message, ip, user_agent)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, contact, name, message, created_at, ip, user_agent
	`

	var c Contact
	err := db.DB.QueryRow(query, contact, name, message, ip, userAgent).Scan(
		&c.ID,
		&c.Contact,
		&c.Name,
		&c.Message,
		&c.CreatedAt,
		&c.IP,
		&c.UserAgent,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *ContactRepository) GetAll() ([]Contact, error) {
	query := `
		SELECT id, contact, name, message, created_at, ip, user_agent
		FROM contacts
		ORDER BY created_at DESC
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var c Contact
		err := rows.Scan(
			&c.ID,
			&c.Contact,
			&c.Name,
			&c.Message,
			&c.CreatedAt,
			&c.IP,
			&c.UserAgent,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

