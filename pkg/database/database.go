package database

import (
	"context"
	"fmt"
	"github.com/Inoi-K/Find-Me/pkg/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *Database

// Database provides ways for interaction with database (currently pool (proxy) only)
type Database struct {
	pool *pgxpool.Pool
}

// ConnectDB creates and returns connection to database
func ConnectDB(ctx context.Context, url string) (*Database, error) {
	if db != nil {
		return db, nil
	}

	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	//defer dbpool.Close()

	db = &Database{
		pool: dbpool,
	}

	return db, nil
}

// GetDB returns database
func GetDB() *Database {
	return db
}

func GetUsers(ctx context.Context) ([]*user.User, error) {
	query := fmt.Sprintf("SELECT id, name FROM public.user;")
	userRows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	users := make([]*user.User, 0, 100)
	for userRows.Next() {
		var userID int
		var name string
		err = userRows.Scan(&userID, &name)
		if err != nil {
			return nil, err
		}

		sphereDescription := make(map[string]string)
		sphereTags := make(map[string]map[string]struct{})

		query = fmt.Sprintf("SELECT sphere_id, description FROM public.user_sphere WHERE user_id = %v;", userID)
		userSphereRows, err := db.pool.Query(ctx, query)
		if err != nil {
			return nil, err
		}
		for userSphereRows.Next() {
			var sphereID int
			var description string
			err = userSphereRows.Scan(&sphereID, &description)
			if err != nil {
				return nil, err
			}

			query = fmt.Sprintf("SELECT name FROM public.sphere WHERE id = %v", sphereID)
			var sphere string
			err = db.pool.QueryRow(ctx, query).Scan(&sphere)
			if err != nil {
				return nil, err
			}
			sphereDescription[sphere] = description

			query = fmt.Sprintf("SELECT name FROM public.tag WHERE id IN (SELECT tag_id FROM public.user_tag WHERE user_id = %v AND sphere_id = %v);", userID, sphereID)
			tagRows, err := db.pool.Query(ctx, query)
			if err != nil {
				return nil, err
			}

			tags := make(map[string]struct{})
			for tagRows.Next() {
				var tag string
				err = tagRows.Scan(&tag)
				if err != nil {
					return nil, err
				}
				tags[tag] = struct{}{}
			}
			sphereTags[sphere] = tags
		}

		usr, err := user.NewUser(name, sphereDescription, sphereTags)
		if err != nil {
			return nil, err
		}
		users = append(users, usr)
	}

	return users, nil
}
