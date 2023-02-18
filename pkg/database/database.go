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

		user, err := model.NewUser(name, sphereDescription, sphereTags)
		if err != nil {
			return nil, err
		}
		users[userID] = user
	}

	return users, nil
}

func AddUser(ctx context.Context, request *pb.SignUpRequest) error {
	query := fmt.Sprintf("INSERT INTO \"user\" VALUES (%d, '%s');", request.UserID, request.Name)
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("INSERT INTO user_sphere VALUES (%d, '%d');", request.UserID, request.SphereID)
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func EditField(ctx context.Context, field, value string, userID, sphereID int64) error {
	query := fmt.Sprintf("UPDATE user_sphere SET %s='%s' WHERE user_id=%d AND sphere_id=%d;", field, value, userID, sphereID)
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func EditTags(ctx context.Context, tags []string, userID, sphereID int64) error {
	var err error
	// TODO goroutine deletion and getting ids
	err = DeleteTags(ctx, userID, sphereID)
	if err != nil {
		return err
	}

	tagIDs, err := ConvertTagsToIDS(ctx, tags)
	if err != nil {
		return err
	}

	return AddTags(ctx, tagIDs, userID, sphereID)
}

func AddTags(ctx context.Context, tags []int64, userID, sphereID int64) error {
	var valuesPart string
	// TODO convert tag names to ids
	for tagID := range tags {
		valuesPart += fmt.Sprintf("(%d, %d, %d),", userID, sphereID, tagID)
	}
	query := fmt.Sprintf("INSERT INTO user_tag VALUES %s;", valuesPart[:len(valuesPart)-1])
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTags(ctx context.Context, userID, sphereID int64) error {
	query := fmt.Sprintf("DELETE FROM user_tag WHERE user_id=%d AND sphere_id=%d", userID, sphereID)
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func ConvertTagsToIDS(ctx context.Context, tags []string) ([]int64, error) {
	var inPart string
	for _, tag := range tags {
		inPart += fmt.Sprintf("'%s',", tag)
	}

	query := fmt.Sprintf("SELECT id FROM tag WHERE name IN (%s);", inPart[:len(inPart)-1])
	tagIDsRows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	// TODO optimize it with caching
	i := 0
	tagIDs := make([]int64, len(tags))
	for tagIDsRows.Next() {
		var tagID int64
		err = tagIDsRows.Scan(&tagID)
		if err != nil {
			return nil, err
		}

		tagIDs[i] = tagID
		i++
	}

	return tagIDs, nil
}
