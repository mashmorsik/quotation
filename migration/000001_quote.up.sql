create table if not exists public.quotation
(
    id uuid primary key,
    base_currency text not null,
    target_currency text not null,
    rate numeric not null,
    time_updated timestamp with time zone not null
);

create table if not exists public.quote_pair
(
    id serial primary key,
    base_currency text not null,
    target_currency text not null
);