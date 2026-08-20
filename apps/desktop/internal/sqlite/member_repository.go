package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(store *Store) *MemberRepository {
	return &MemberRepository{db: store.db}
}

func (r *MemberRepository) Create(ctx context.Context, member domain.Member) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO members (
  id, gym_id, first_name, last_name, email, phone, identification_number,
  date_of_birth, address, status, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, member.ID(), member.GymID(), member.FirstName(), member.LastName(), nullIfEmpty(member.Email()), member.Phone(), nullIfEmpty(member.IdentificationNumber()), nullIfEmpty(member.DateOfBirth()), nullIfEmpty(member.Address()), member.Status(), domain.FormatTimestamp(member.CreatedAt()), domain.FormatTimestamp(member.UpdatedAt()))
	if err != nil {
		return fmt.Errorf("insert member: %w", err)
	}
	return nil
}

func (r *MemberRepository) Get(ctx context.Context, gymID domain.GymID, memberID domain.MemberID) (domain.Member, error) {
	member, err := scanMember(r.db.QueryRowContext(ctx, memberSelect+` WHERE gym_id = ? AND id = ? AND deleted_at IS NULL`, gymID, memberID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Member{}, ports.ErrMemberNotFound
	}
	if err != nil {
		return domain.Member{}, fmt.Errorf("get member: %w", err)
	}
	return member, nil
}

func (r *MemberRepository) List(ctx context.Context, gymID domain.GymID) ([]domain.Member, error) {
	rows, err := r.db.QueryContext(ctx, memberSelect+` WHERE gym_id = ? AND deleted_at IS NULL ORDER BY created_at DESC`, gymID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()

	members := []domain.Member{}
	for rows.Next() {
		member, err := scanMember(rows)
		if err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate members: %w", err)
	}
	return members, nil
}

func (r *MemberRepository) Update(ctx context.Context, member domain.Member) error {
	result, err := r.db.ExecContext(ctx, `
UPDATE members
SET first_name = ?, last_name = ?, email = ?, phone = ?, identification_number = ?,
    date_of_birth = ?, address = ?, status = ?, updated_at = ?, version = version + 1
WHERE gym_id = ? AND id = ? AND deleted_at IS NULL
`, member.FirstName(), member.LastName(), nullIfEmpty(member.Email()), member.Phone(), nullIfEmpty(member.IdentificationNumber()), nullIfEmpty(member.DateOfBirth()), nullIfEmpty(member.Address()), member.Status(), domain.FormatTimestamp(member.UpdatedAt()), member.GymID(), member.ID())
	if err != nil {
		return fmt.Errorf("update member: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("update member rows affected: %w", err)
	} else if affected == 0 {
		return ports.ErrMemberNotFound
	}
	return nil
}

func (r *MemberRepository) Archive(ctx context.Context, gymID domain.GymID, memberID domain.MemberID, archivedAt time.Time) error {
	archivedAt = archivedAt.UTC()
	result, err := r.db.ExecContext(ctx, `
UPDATE members
SET deleted_at = ?, updated_at = ?, version = version + 1
WHERE gym_id = ? AND id = ? AND deleted_at IS NULL
`, domain.FormatTimestamp(archivedAt), domain.FormatTimestamp(archivedAt), gymID, memberID)
	if err != nil {
		return fmt.Errorf("archive member: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("archive member rows affected: %w", err)
	} else if affected == 0 {
		return ports.ErrMemberNotFound
	}
	return nil
}

const memberSelect = `
SELECT id, gym_id, first_name, last_name, email, phone, identification_number,
       date_of_birth, address, status, created_at, updated_at
FROM members`

type rowScanner interface {
	Scan(...any) error
}

func scanMember(row rowScanner) (domain.Member, error) {
	var id, gymID, firstName, lastName, phone, status, createdAt, updatedAt string
	var email, identificationNumber, dateOfBirth, address sql.NullString
	if err := row.Scan(&id, &gymID, &firstName, &lastName, &email, &phone, &identificationNumber, &dateOfBirth, &address, &status, &createdAt, &updatedAt); err != nil {
		return domain.Member{}, err
	}
	memberID, err := domain.ParseMemberID(id)
	if err != nil {
		return domain.Member{}, fmt.Errorf("parse member ID: %w", err)
	}
	ownerID, err := domain.ParseGymID(gymID)
	if err != nil {
		return domain.Member{}, fmt.Errorf("parse gym ID: %w", err)
	}
	created, err := domain.ParseTimestamp(createdAt)
	if err != nil {
		return domain.Member{}, err
	}
	updated, err := domain.ParseTimestamp(updatedAt)
	if err != nil {
		return domain.Member{}, err
	}
	return domain.NewMember(memberID, ownerID, firstName, lastName, email.String, phone, identificationNumber.String, dateOfBirth.String, address.String, domain.MemberStatus(status), created, updated)
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
