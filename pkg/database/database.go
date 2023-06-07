package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func GetWeights(ctx context.Context) (map[int64]map[int64]float64, error) {
	query := fmt.Sprintf("SELECT * FROM weight;")
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	//defer rows.Close()

	w := make(map[int64]map[int64]float64)
	for rows.Next() {
		var fromID, toID int64
		var weight float64
		err = rows.Scan(&fromID, &toID, &weight)
		if err != nil {
			return w, err
		}

		if _, ok := w[fromID]; !ok {
			w[fromID] = make(map[int64]float64)
		}
		w[fromID][toID] = weight
	}
	rows.Close()

	query = fmt.Sprintf("SELECT id, weight FROM dimension;")
	rows, err = db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dimensionID int64
		var weight float64
		err = rows.Scan(&dimensionID, &weight)
		if err != nil {
			return w, err
		}

		w[dimensionID] = map[int64]float64{0: weight}
	}

	return w, nil
}

// GetSearchFamiliar gets a search option for recommendations generation
func GetSearchFamiliar(ctx context.Context, userID, sphereID int64) (searchFamiliar bool, err error) {
	// TODO make it categorical field 'search_option'
	query := fmt.Sprintf("SELECT search_familiar FROM user_sphere WHERE user_id=%d AND sphere_id=%d;", userID, sphereID)
	err = db.pool.QueryRow(ctx, query).Scan(&searchFamiliar)
	return
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

// GetUsersTag returns user-sphere-tag model of all users
func GetUsersTag(ctx context.Context) (model.USDT, error) {
	query := fmt.Sprintf("SELECT * FROM user_tag;")
	return getUserTag(ctx, query)
}

// GetUserTag returns user-sphere-tag model of users for the given one
func GetUserTag(ctx context.Context, userID, sphereID int64) (model.USDT, error) {
	query := fmt.Sprintf("SELECT * FROM user_tag WHERE user_id=%d AND sphere_id=%d;", userID, sphereID)
	return getUserTag(ctx, query)
}

// getUserTag converts select query to user-sphere-tag model
func getUserTag(ctx context.Context, query string) (model.USDT, error) {
	userSphereTagRows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer userSphereTagRows.Close()

	usdt := make(model.USDT, 1000)
	for userSphereTagRows.Next() {
		var userID, sphereID, tagID, dimensionID int64
		err = userSphereTagRows.Scan(&userID, &sphereID, &tagID, &dimensionID)
		if err != nil {
			return nil, err
		}

		// new user
		if _, ok := usdt[userID]; !ok {
			usdt[userID] = make(map[int64]map[int64]map[int64]struct{}, 1000)
		}
		// new sphere
		if _, ok := usdt[userID][sphereID]; !ok {
			usdt[userID][sphereID] = make(map[int64]map[int64]struct{})
		}
		// new dimension
		if _, ok := usdt[userID][sphereID][dimensionID]; !ok {
			usdt[userID][sphereID][dimensionID] = make(map[int64]struct{}, config.C.TagsLimit)
		}
		usdt[userID][sphereID][dimensionID][tagID] = struct{}{}
	}

	return usdt, nil
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
		query = fmt.Sprintf("SELECT name FROM tag WHERE id IN (SELECT tag_id FROM user_tag WHERE user_id=%d AND sphere_id=%d AND dimension_id=%d);", userID, sphereID, 1)
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
	query := fmt.Sprintf("INSERT INTO \"user\" VALUES (%d, '%s', '%s', %s, '%s', '%s', '%s');", request.UserID, request.Name, request.Gender, request.Age, request.Faculty, request.University, request.Username)
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

// AddTags creates new tags for a user in a sphere
func AddTags(ctx context.Context, tags []string, userID, sphereID, dimensionID int64) error {
	var valuesPart string
	for _, tagID := range tags {
		valuesPart += fmt.Sprintf("(%d, %d, %s, %d),", userID, sphereID, tagID, dimensionID)
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

// GetAssociateTags return associated to the given tags
func GetAssociateTags(ctx context.Context, tags []string) ([]string, error) {
	t := strings.Join(tags, ",")
	query := fmt.Sprintf("SELECT associated_id FROM association_rule WHERE base_id IN (%s) AND associated_id NOT IN (%s);", t, t)
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	addTagIDs := make([]string, 0)
	for rows.Next() {
		var addTagID string
		err = rows.Scan(&addTagID)
		if err != nil {
			return addTagIDs, err
		}
		addTagIDs = append(addTagIDs, addTagID)
	}

	return addTagIDs, nil
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

func GetMatches(ctx context.Context, sphereID int64) (map[int64]map[int64]bool, error) {
	query := fmt.Sprintf("SELECT from_id, to_id, is_like FROM match WHERE sphere_id=%d;", sphereID)
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make(map[int64]map[int64]bool)
	for rows.Next() {
		var fromID, toID int64
		var isLike bool
		err = rows.Scan(&fromID, &toID, &isLike)
		if err != nil {
			return nil, err
		}

		if _, ok := matches[fromID]; !ok {
			matches[fromID] = make(map[int64]bool)
		}

		matches[fromID][toID] = isLike
	}

	return matches, nil
}

func GetUserMatches(ctx context.Context, fromID, sphereID int64) (map[int64]bool, error) {
	query := fmt.Sprintf("SELECT to_id, is_like FROM match WHERE from_id=%d AND sphere_id=%d;", fromID, sphereID)
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make(map[int64]bool)
	for rows.Next() {
		var toID int64
		var isLike bool
		err = rows.Scan(&toID, &isLike)
		if err != nil {
			return nil, err
		}

		matches[toID] = isLike
	}

	return matches, nil
}
