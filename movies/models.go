package movies

import "time"

type Movie struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Director  []string    `json:"director"`
	Year      string    `json:"year"`
	Userid    string    `json:"userid"`
	CreatedOn time.Time `json:"createdon"`
	UpdatedOn time.Time `json:"updatedon"`
}

/*
CREATE TABLE movies (
id UUID NOT NULL DEFAULT gen_random_uuid(),
title STRING NOT NULL,
year STRING NOT NULL,
userid UUID NOT NULL,
createdon TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
updatedon TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
CONSTRAINT "primary" PRIMARY KEY (id ASC))
*/

/*
CREATE TABLE movie_directors (
movie_id UUID NOT NULL,
director STRING NOT NULL,
createdon TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
updatedon TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now())
*/