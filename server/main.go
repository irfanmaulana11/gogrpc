package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogrpc/config"
	"github.com/gogrpc/model"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type server struct {
	conn *sql.DB
}

func (connection *server) GetUserProfile(ctx context.Context, req *model.GetUserProfileRequest) (*model.UserProfile, error) {
	db := connection.conn
	id := req.GetId()
	sqlStatement := `SELECT id, First_name, last_name, email, password FROM "profile" WHERE "id" = $1`

	var profile model.UserProfile
	statement, err := db.Prepare(sqlStatement)

	if err != nil {
		return nil, err
	}
	//defer statement.Close()
	err = statement.QueryRow(id).Scan(&profile.ID, &profile.FirstName, &profile.LastName, &profile.Email, &profile.Password)

	if err != nil {
		return nil, err
	}
	res := &model.UserProfile{
		//ID:        profile.ID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Email:     profile.Email,
		Password:  profile.Password,
		//Createdat: profile.Createdat,
		//Updatedat: profile.Updatedat,
	}
	return res, nil
}

func (connection *server) CreateUserProfile(ctx context.Context, req *model.CreateUserProfileRequest) (*model.UserProfile, error) {
	db := connection.conn
	firstname := req.GetUserProfile().GetFirstName()
	lastname := req.GetUserProfile().GetLastName()
	email := req.GetUserProfile().GetEmail()
	password := req.GetUserProfile().GetPassword()
	sqlStatement := `INSERT INTO "profile" ( first_name, last_name, email, password) VALUES ($1, $2, $3, $4)`
	if _, err := db.Exec(sqlStatement, firstname, lastname, email, password); err != nil {
		return nil, errors.Wrap(err, "User couldn't be inserted")
	}
	return req.UserProfile, nil
}

func (connection *server) DeleteUserProfile(ctx context.Context, req *model.DeleteUserProfileRequest) (*empty.Empty, error) {
	db := connection.conn
	id := req.GetId()
	sqlStatement := `delete from "profile" where id=$1`
	if _, err := db.Exec(sqlStatement, id); err != nil {
		errors.Wrap(err, "User couldn't be deleted")
	}
	return &empty.Empty{}, nil
}

func (connection *server) UpdateUserProfile(ctx context.Context, req *model.UpdateUserProfileRequest) (*model.UserProfile, error) {
	db := connection.conn
	sqlStatement := `UPDATE "profile" SET first_name=$1, last_name=$2, email=$3, password=$4 WHERE "id" =$5;`
	if _, err := db.Exec(sqlStatement, req.UserProfile.FirstName, req.UserProfile.LastName, req.UserProfile.Email, req.UserProfile.Password, req.GetUserProfile().GetID()); err != nil {
		return nil, err
	}
	return req.UserProfile, nil
}

func (connection *server) ListUsersProfiles(ctx context.Context, req *model.ListUsersProfilesRequest) (*model.ListUsersProfilesResponse, error) {
	db := connection.conn
	//id := req.GetQuery() + "%"
	sqlStatement := ` select first_name, last_name, email, password from "profile"`
	result, err := db.Query(sqlStatement)
	defer result.Close()
	if err != nil {
		fmt.Println(err)
	}
	res := []*model.UserProfile{}
	for result.Next() {
		var first, last, email, id string
		if err = result.Scan(&first, &last, &email, &id); err != nil {
			errors.Wrap(err, "Users couln't be listed")
		}
		u := model.UserProfile{
			FirstName: first,
			LastName:  last,
			Email:     email,
			ID:        id,
		}
		res = append(res, &u)
	}
	ans := model.ListUsersProfilesResponse{
		Profiles: res,
	}
	return &ans, nil
}

func main() {
	fmt.Println("Welcome to the server")
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		errors.Wrap(err, " Failed to listen the port")
	}
	s := grpc.NewServer()
	db, err := config.GetConnection()

	if err != nil {
		errors.Wrap(err, "Connection couldn't be opened")
	} else {
		fmt.Println("Connected to DB")
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		errors.Wrap(err, "Connection not established, ping didn't work")
	}
	model.RegisterUserProfilesServer(s, &server{db})
	if err := s.Serve(lis); err != nil {
		errors.Wrap(err, "Failed to server the listener")
	}
}
