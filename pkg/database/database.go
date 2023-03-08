package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *Database

// Database provides ways for interaction with database (currently pool (proxy) only)
type Database struct {
	pool *pgxpool.Pool
}

// ConnectDB creates connection to database
func ConnectDB(ctx context.Context) error {
	if db != nil {
		return nil
	}

	dbConfig, err := pgxpool.ParseConfig(config.C.DatabaseURL)
	if err != nil {
		return err
	}
	dbConfig.MaxConns = config.C.DatabasePoolMaxConnections
	dbpool, err := pgxpool.NewWithConfig(ctx, dbConfig)
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
	query := fmt.Sprintf("SELECT 1 FROM \"user\" WHERE id=%d;", userID)
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// GetUserSphereTag returns user-sphere-tag model of all users
func GetUserSphereTag(ctx context.Context) (model.UST, error) {
	query := fmt.Sprintf("SELECT * FROM user_tag;")
	return getUserSphereTag(ctx, query)
}

// GetNewUserSphereTag returns user-sphere-tag model of new (no record in match table) users for the given one
func GetNewUserSphereTag(ctx context.Context, userID, sphereID int64) (model.UST, error) {
	query := fmt.Sprintf("SELECT * FROM user_tag WHERE user_id NOT IN (SELECT to_id FROM match WHERE from_id=%d AND sphere_id=%d);", userID, sphereID)
	return getUserSphereTag(ctx, query)
}

// getUserSphereTag converts select query to user-sphere-tag model
func getUserSphereTag(ctx context.Context, query string) (model.UST, error) {
	userSphereTagRows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer userSphereTagRows.Close()

	ust := make(model.UST, 1000)
	for userSphereTagRows.Next() {
		var userID, sphereID, tagID int64
		err = userSphereTagRows.Scan(&userID, &sphereID, &tagID)
		if err != nil {
			return nil, err
		}

		// new user
		if _, ok := ust[userID]; !ok {
			ust[userID] = make(map[int64]map[int64]struct{}, 1000)
		}
		// new sphere
		if _, ok := ust[userID][sphereID]; !ok {
			ust[userID][sphereID] = make(map[int64]struct{}, config.C.TagsLimit)
		}
		ust[userID][sphereID][tagID] = struct{}{}
	}

	return ust, nil
}

// GetUsers gets all (main and additional) information about all user
// TODO caching/sharding/partitioning?
func GetUsers(ctx context.Context) (map[int64]*model.User, error) {
	query := fmt.Sprintf("SELECT id FROM \"user\";")
	userRows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	// TODO sus 100
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
	query := fmt.Sprintf("SELECT name, gender, age, faculty, university, username FROM \"user\" WHERE id=%d;", userID)
	err := db.pool.QueryRow(ctx, query).Scan(&user.Name, &user.Gender, &user.Age, &user.Faculty, &user.University, &user.Username)
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
	defer userSphereRows.Close()

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

		// get tags
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
		tagRows.Close()

		// combine fields in a model
		sphereInfo[sphereID] = &model.UserSphere{
			Description: description,
			PhotoID:     photo,
			Tags:        tags,
		}
	}
	return sphereInfo, nil
}

// AddUser adds user to db in all necessary tables
// TODO refactor pb req to user model
func AddUser(ctx context.Context, request *pb.SignUpRequest) error {
	query := fmt.Sprintf("INSERT INTO \"user\" VALUES (%d, '%s', '%s', %s, '%s');", request.UserID, request.Name, request.Gender, request.Age, request.Faculty)
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
		query = fmt.Sprintf("UPDATE \"user\" SET %s='%s' WHERE id=%d;", field, value, userID)
		_, err = db.pool.Query(ctx, query)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetTags returns tags associated to the sphere
func GetTags(ctx context.Context, sphereID int64) ([]model.Tag, error) {
	query := fmt.Sprintf("SELECT * FROM tag WHERE id IN (SELECT tag_id FROM sphere_tag WHERE sphere_id=%d);", sphereID)
	tagRows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	tags := make([]model.Tag, 0, 100)
	for tagRows.Next() {
		var id string
		var name string
		err = tagRows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		tag := model.Tag{
			ID:   id,
			Name: name,
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// EditTags updates tags of a user
func EditTags(ctx context.Context, tagIDs []string, userID, sphereID int64) error {
	var err error
	err = DeleteTags(ctx, userID, sphereID)
	if err != nil {
		return err
	}

	return AddTags(ctx, tagIDs, userID, sphereID)
}

// AddTags creates new tags for a user in a sphere
func AddTags(ctx context.Context, tags []string, userID, sphereID int64) error {
	var valuesPart string
	for _, tagID := range tags {
		valuesPart += fmt.Sprintf("(%d, %d, %s),", userID, sphereID, tagID)
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
	query := fmt.Sprintf("DELETE FROM user_tag WHERE user_id=%d AND sphere_id=%d;", userID, sphereID)
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
	defer tagIDsRows.Close()

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

// Match registers like in the match table and checks for the other way like
func Match(ctx context.Context, fromID, toID, sphereID int64, isLike bool) (bool, error) {
	// remove the previous record if exists
	query := fmt.Sprintf("DELETE FROM match WHERE from_id=%d AND to_id=%d;", fromID, toID)
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return false, err
	}

	// record a new like
	query = fmt.Sprintf("INSERT INTO match VALUES (%d, %d, %t, %d);", fromID, toID, isLike, sphereID)
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return false, err
	}

	// check if there is a return like
	query = fmt.Sprintf("SELECT is_like FROM match WHERE from_id=%d AND to_id=%d;", toID, fromID)
	isReciprocated := false
	err = db.pool.QueryRow(ctx, query).Scan(&isReciprocated)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}

	return isReciprocated, nil
}
