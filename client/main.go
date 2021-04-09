package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gogrpc/model"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("welcome to Client app")
	conn, err := grpc.Dial("localhost:50005", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	client := model.NewUserProfilesClient(conn)

	g := gin.Default()

	g.GET("/user/get", func(ctx *gin.Context) {
		req := &model.ListUsersProfilesRequest{}
		res, err := client.ListUsersProfiles(context.Background(), req)

		if err != nil {
			log.Fatalf("Error while Geting user profile %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, res)
	})

	g.GET("/user/get/:id", func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		//idUser, _ := strconv.ParseInt(idParam, 10, 64)
		req := &model.GetUserProfileRequest{
			Id: idParam,
		}
		res, err := client.GetUserProfile(context.Background(), req)

		if err != nil {
			log.Fatalf("Error while getting user profile %v", err)
		}
		//fmt.Println(res)
		ctx.JSON(http.StatusOK, res)

	})

	g.POST("/user/post", func(ctx *gin.Context) {
		var user model.UserProfile
		bindErr := ctx.BindJSON(&user)
		if bindErr != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprint(bindErr)) // jika tidak sesuai struct maka status 400
			return
		}

		req := &model.CreateUserProfileRequest{
			UserProfile: &model.UserProfile{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Password:  user.Password,
			},
		}

		res, err := client.CreateUserProfile(context.Background(), req)

		if err != nil {
			log.Fatalf("Error while creating user profile %v", err)
		}
		//fmt.Println(res)
		ctx.JSON(http.StatusOK, res)
	})

	g.PUT("user/put/:id", func(ctx *gin.Context) {
		var user model.UserProfile
		idParam := ctx.Param("id")

		bindErr := ctx.BindJSON(&user)
		if bindErr != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprint(bindErr))
			return
		}

		req := &model.UpdateUserProfileRequest{
			UserProfile: &model.UserProfile{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Password:  user.Password,
				ID:        idParam,
			},
		}
		res, err := client.UpdateUserProfile(context.Background(), req)

		if err != nil {
			log.Fatalf("Error while Updating user profile %v", err)
		}
		ctx.JSON(http.StatusOK, res)
	})

	g.DELETE("user/delete/:id", func(ctx *gin.Context) {
		idParam := ctx.Param("id")

		req := &model.DeleteUserProfileRequest{Id: idParam}
		res, err := client.DeleteUserProfile(context.Background(), req)
		if err != nil {
			log.Fatalf("Error while Geting user profile %v", err)
		}
		ctx.JSON(http.StatusOK, res)
	})

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to Run server: %v", err)
	}
}
