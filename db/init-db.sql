create table Users (
    userId bigint generated by default as identity not null,
    externalId varchar(100) not null, 
    constraint PK_Users primary key (userId)
);

create index idx_Users_externalId on Users (externalId);

create table ShortUrls (

    shortUrlId bigint generated by default as identity not null,
    userId bigint not null, 
    fullUrl varchar(2000) not null,
    shortUrl varchar(50) not null,
    constraint PK_ShortUrls primary key (shortUrlId)
);

create unique index uniq_ShortUrls_shortUrl on ShortUrls (shortUrl);

create index idx_ShortUrls_userId on ShortUrls (userId);