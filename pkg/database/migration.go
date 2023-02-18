package database

import "context"

func CreateTables(ctx context.Context) error {
	// user
	query := "create table if not exists \"user\"\n(\n    id   integer not null\n        constraint user_pk\n            primary key,\n    name text    not null\n);"
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	// sphere
	query = "create table if not exists sphere\n(\n    id   serial\n        constraint sphere_pk\n            primary key,\n    name text not null\n);"
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	// tag
	query = "create table if not exists tag\n(\n    id   serial\n        constraint tag_pk\n            primary key,\n    name text not null\n);"
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	// user_sphere
	query = "create table if not exists user_sphere\n(\n    user_id     integer not null\n        constraint user_sphere_user_fk\n            references \"user\",\n    sphere_id   integer not null\n        constraint user_sphere_sphere_fk\n            references sphere\n            on update cascade on delete cascade,\n    description text,\n    photo       text,\n    constraint user_sphere_pk\n        primary key (user_id, sphere_id)\n);"
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	// sphere_tag
	query = "create table if not exists sphere_tag\n(\n    sphere_id integer not null\n        constraint sphere_tag_sphere_fk\n            references sphere\n            on update cascade on delete cascade,\n    tag_id    integer not null\n        constraint sphere_tag_tag_fk\n            references tag\n            on update cascade on delete cascade,\n    constraint sphere_tag_pk\n        primary key (sphere_id, tag_id)\n);"
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	// user_tag
	query = "create table if not exists user_tag\n(\n    user_id   integer not null\n        constraint user_sphere_tag_user_fk\n            references \"user\"\n            on update cascade on delete cascade,\n    sphere_id integer not null\n        constraint user_sphere_tag_sphere_fk\n            references sphere\n            on update cascade on delete cascade,\n    tag_id    integer not null\n        constraint user_sphere_tag_tag_fk\n            references tag\n            on update cascade on delete cascade,\n    constraint user_sphere_tag_pk\n        primary key (user_id, sphere_id, tag_id)\n);"
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	// match
	query = "create table if not exists match\n(\n    user_id_1 integer not null\n        constraint match_user_id_1_fk\n            references \"user\"\n            on update cascade on delete cascade,\n    user_id_2 integer not null\n        constraint match_user_id_2_fk\n            references \"user\"\n            on update cascade on delete cascade,\n    constraint match_pk\n        primary key (user_id_1, user_id_2)\n);"
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
