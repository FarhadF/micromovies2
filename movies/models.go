package movies

import (
	"time"
)

type Movie struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Director  []string    `json:"director"`
	Year      string    `json:"year"`
	CreatedBy    string    `json:"createdby"`
	CreatedOn time.Time `json:"createdon"`
	UpdatedBy 	string  `json:"updatedby"`
	UpdatedOn time.Time `json:"updatedon"`
}

/*
CREATE TABLE movies (
id UUID NOT NULL DEFAULT gen_random_uuid(),
title STRING NOT NULL,
year STRING NOT NULL,
createdby UUID NOT NULL,
createdon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
updatedby UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
updatedon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT '1970-01-01',
CONSTRAINT "primary" PRIMARY KEY (id ASC))
*/

/*
CREATE TABLE movie_directors (
movie_id UUID NOT NULL,
director STRING NOT NULL,
createdby UUID NOT NULL,
createdon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now())
*/