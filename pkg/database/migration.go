package database

import "context"

// CreateTables creates base tables
func CreateTables(ctx context.Context) error {
	queries := []string{
		// user
		"create table if not exists \"user\"\n(\n    id         integer    not null\n        constraint user_pk\n            primary key,\n    name       text       not null,\n    gender     varchar(1) not null,\n    age        integer    not null,\n    faculty    text       not null,\n    university text       not null,\n    username   text       not null\n);",
		// sphere
		"create table if not exists sphere\n(\n    id   serial\n        constraint sphere_pk\n            primary key,\n    name text not null\n);",
		// tag
		"create table if not exists tag\n(\n    id   serial\n        constraint tag_pk\n            primary key,\n    name text not null\n);",
		// dimension
		"create table if not exists dimension\n(\n    id     serial\n        constraint dimension_pk\n            primary key,\n    name   text             not null,\n    weight double precision not null\n);",
		// association rule
		"create table if not exists association_rule\n(\n    base_id       integer not null\n        constraint table_name_tag_id_fk\n            references tag,\n    associated_id integer not null\n        constraint table_name_tag_id_fk_2\n            references tag,\n    constraint association_rule_pk\n        primary key (base_id, associated_id)\n);",
		// user_sphere
		"create table if not exists user_sphere\n(\n    user_id     integer not null\n        constraint user_sphere_user_fk\n            references \"user\",\n    sphere_id   integer not null\n        constraint user_sphere_sphere_fk\n            references sphere\n            on update cascade on delete cascade,\n    description text,\n    photo       text,\n    constraint user_sphere_pk\n        primary key (user_id, sphere_id)\n);",
		// sphere_tag
		"create table if not exists sphere_tag\n(\n    sphere_id integer not null\n        constraint sphere_tag_sphere_fk\n            references sphere\n            on update cascade on delete cascade,\n    tag_id    integer not null\n        constraint sphere_tag_tag_fk\n            references tag\n            on update cascade on delete cascade,\n    constraint sphere_tag_pk\n        primary key (sphere_id, tag_id)\n);",
		// user_tag
		"create table if not exists user_tag\n(\n    user_id      integer not null\n        constraint user_sphere_tag_user_fk\n            references \"user\"\n            on update cascade on delete cascade,\n    sphere_id    integer not null\n        constraint user_sphere_tag_sphere_fk\n            references sphere\n            on update cascade on delete cascade,\n    tag_id       integer not null\n        constraint user_sphere_tag_tag_fk\n            references tag\n            on update cascade on delete cascade,\n    dimension_id integer not null\n        constraint user_tag_dimension_id_fk\n            references dimension\n            on update cascade on delete cascade,\n    constraint user_tag_pk\n        primary key (user_id, sphere_id, tag_id, dimension_id)\n);",
		// match
		"create table if not exists match\n(\n    from_id   integer not null\n        constraint match_user_id_1_fk\n            references \"user\"\n            on update cascade on delete cascade,\n    to_id     integer not null\n        constraint match_user_id_2_fk\n            references \"user\"\n            on update cascade on delete cascade,\n    is_like   boolean not null,\n    sphere_id integer not null\n        constraint match_sphere_id_fk\n            references sphere,\n    constraint match_pk\n        primary key (from_id, to_id)\n);",
	}
	for _, query := range queries {
		_, err := db.pool.Query(ctx, query)
		if err != nil {
			return err
		}
	}

	return nil
}
