package users

import "time"

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedOn time.Time `json:"createdon"`
	UpdatedOn time.Time `json:"updatedon"`
	LastLogin time.Time `json:"lastlogin"`
}

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
CREATE TABLE users (
id UUID NOT NULL DEFAULT gen_random_uuid(),
name STRING NULL,
lastname STRING NULL,
email STRING NOT NULL,
password STRING NOT NULL,
userrole string NOT NULL,
createdon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
updatedon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT '1970-01-01',
lastlogin TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT '1970-01-01',
CONSTRAINT "primary" PRIMARY KEY (id ASC))
*/
