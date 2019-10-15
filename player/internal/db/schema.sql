create table rounds (
    id bigint not null auto_increment,
    external_id bigint not null,
    player varchar (255),
    `status` int not null,
    created_at datetime not null,
    updated_at datetime not null,
    
    primary key(id),
);

create table round_events (
    id bigint not null auto_increment,
    foreign_id bigint not null,
    `type` int not null,
    updated_at datetime not null,

    primary key(id)
);

create table parts (
    id bigint not null auto_increment,
    round_id bigint not null,
    player varchar(255) not null,
    value int not null,
    rank int,
    submitted bool not null,
    created_at datetime not null,
    updated_at datetime not null,

    primary key(id)
);

create table cursors (
    id varchar(255) not null,
    last_event_id bigint not null,
    updated_at datetime not null,

    primary key(id)
);