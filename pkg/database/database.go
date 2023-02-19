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

// UserExists checks if user exists in db
func UserExists(ctx context.Context, userID int64) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM user WHERE id=%d;", userID)
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

// GetUsers gets all (main and additional) information about all user
func GetUsers(ctx context.Context) (map[int64]*model.User, error) {
	query := fmt.Sprintf("SELECT id FROM user;")
	userRows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	users := make(map[int64]*model.User, 100)
	for userRows.Next() {
		var userID int64
		err = userRows.Scan(&userID)
		if err != nil {
			return nil, err
		}

		users[userID], err = GetUser(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

// GetUser gets all (main and additional) information about user
func GetUser(ctx context.Context, userID int64) (*model.User, error) {
	user, err := GetUserMain(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.SphereInfo, err = GetUserAdditional(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserMain gets main info about a user
func GetUserMain(ctx context.Context, userID int64) (*model.User, error) {
	user := &model.User{}
	query := fmt.Sprintf("SELECT name, gender, age, faculty FROM \"user\" WHERE id=%d", userID)
	err := db.pool.QueryRow(ctx, query).Scan(&user.Name, &user.Gender, &user.Age, &user.Faculty)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserAdditional gets additional info of a user in a sphere
func GetUserAdditional(ctx context.Context, userID int64) (map[int64]*model.UserSphere, error) {
	sphereInfo := make(map[int64]*model.UserSphere)
	query := fmt.Sprintf("SELECT sphere_id, description, photo FROM user_sphere WHERE user_id=%d;", userID)
	userSphereRows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	// TODO optimize it with caching
	for userSphereRows.Next() {
		var (
			sphereID    int64
			description string
			photo       string
		)
		err = userSphereRows.Scan(&sphereID, &description, &photo)
		if err != nil {
			return nil, err
		}

		query = fmt.Sprintf("SELECT name FROM tag WHERE id IN (SELECT tag_id FROM user_tag WHERE user_id=%d AND sphere_id=%d);", userID, sphereID)
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

		sphereInfo[sphereID] = &model.UserSphere{
			Description: description,
			PhotoID:     photo,
			Tags:        tags,
		}
	}
	return sphereInfo, nil
}

// AddUser adds user to db in all necessary tables
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

// EditField updates additional text-based value of user in a sphere in db
func EditField(ctx context.Context, field, value string, userID, sphereID int64) error {
	query := fmt.Sprintf("UPDATE user_sphere SET %s='%s' WHERE user_id=%d AND sphere_id=%d;", field, value, userID, sphereID)
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// EditTags updates tags of a user
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

// AddTags creates new tags for a user in a sphere
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

// DeleteTags deletes previous tags of a user
func DeleteTags(ctx context.Context, userID, sphereID int64) error {
	query := fmt.Sprintf("DELETE FROM user_tag WHERE user_id=%d AND sphere_id=%d", userID, sphereID)
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

// ConvertTagsToIDS gets names of tags and returns their IDs
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
