package database

import (
	"context"
	"fmt"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *Database

// Database provides ways for interaction with database (currently pool (proxy) only)
type Database struct {
	pool *pgxpool.Pool
}

// ConnectDB creates connection to database
func ConnectDB(ctx context.Context, url string) error {
	if db != nil {
		return nil
	}

	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		return err
	}
	//defer dbpool.Close()

	db = &Database{
		pool: dbpool,
	}

	return CreateTables(ctx)
}

func UserExists(ctx context.Context, userID int64) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM public.user WHERE id=%d;", userID)
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func GetUsers(ctx context.Context) (map[int64]*model.User, error) {
	query := fmt.Sprintf("SELECT id, name FROM public.user;")
	userRows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	users := make(map[int64]*model.User, 100)
	for userRows.Next() {
		var userID int64
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
		// TODO optimize it with caching
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

		usr, err := model.NewUser(name, sphereDescription, sphereTags)
		if err != nil {
			return nil, err
		}
		users[userID] = usr
	}

	return users, nil
}

func AddUser(ctx context.Context, request *pb.SignUpRequest) error {
	query := fmt.Sprintf("INSERT INTO \"user\" VALUES (%d, '%s');", request.UserID, request.Name)
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
