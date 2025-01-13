package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// Define a Permissions slice, which we will use to hold the permission codes (like
// "movies:read" and "movies:write") for a single user.
type Permissions []string

// Define the PermissionModel type.
type PermissionModel struct {
	DB *sql.DB
}

// Add a helper method to check whether the Permissions slice contains a specific
// permission code.
func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

// The GetAllForUser() method returns all permission codes for a specific user in a
// Permissions slice. The code in this method should feel very familiar --- it uses the
// standard pattern that we've already seen before for retrieving multiple data rows in
// an SQL query.
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
		SELECT p.code
		FROM users u
		INNER JOIN users_permissions up ON u.id = up.user_id
		INNER JOIN permissions p ON up.permission_id = p.id
		WHERE u.id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := Permissions{}

	for rows.Next() {
		var permission string 
		err = rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	query := `
		INSERT INTO users_permissions
		SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)
	`

	args := []any{userID, pq.Array(codes)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}