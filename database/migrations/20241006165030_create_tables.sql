-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id serial primary key,
    name varchar(20)
);

create table if not exists mood_scores (
    id serial primary key,
    score smallint,
    date date not null default current_date,
    user_id integer,
    foreign key (user_id) references users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists mood_scores;
drop table if exists users;
-- +goose StatementEnd
