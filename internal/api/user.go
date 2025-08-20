package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/mauzec/user-api/db/sqlc"
	"github.com/mauzec/user-api/internal/token"
	"github.com/mauzec/user-api/internal/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32,alphanum"`
	Fullname string `json:"fullname" binding:"required,min=3,max=64"`
	Sex      string `json:"sex" binding:"required,sex"`
	Age      int32  `json:"age" binding:"required,min=18,max=60"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,phone"`
	Password string `json:"password" binding:"required,min=5,max=64"`
}

type userResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullname"`
	Sex       string    `json:"sex"`
	Age       int32     `json:"age"`
	Avatar    string    `json:"avatar"`
	Status    string    `json:"status"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Sex:       user.Sex,
		Age:       user.Age,
		Avatar:    user.Avatar,
		Status:    user.Status,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Time,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid request body")))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("something went wrong")))
		return
	}

	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.Fullname,
		Sex:            req.Sex,
		Age:            req.Age,
		Email:          req.Email,
		Phone:          req.Phone,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("this user already exists")))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("something went wrong")))
		return
	}

	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)
}

type loginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidRequest))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrNotFound))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServerError))
		return
	}

	err = util.CheckPassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid username or password")))
		return
	}

	token, err := server.tokenMaker.CreateToken(
		req.Username,
		server.tokenParams.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServerError))
		return
	}

	resp := loginResponse{
		Token: token,
		User:  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, resp)
}

type getUserUri struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) getUserByUsername(ctx *gin.Context) {
	var uri getUserUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrInvalidRequest)
		return
	}

	payloadVal, exists := ctx.Get(authPayloadKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrMissingAuthPayload))
		return
	}
	_, ok := payloadVal.(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrMissingAuthPayload))
		return
	}

	// use if u want a user to see only his own data
	// if payload.Username != uri.Username {
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("permission denied")))
	// 	return
	// }

	user, err := server.store.GetUserByUsername(ctx, uri.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type updateUserRequest struct {
	Fullname *string `json:"fullname" binding:"omitempty,min=3,max=64"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone" binding:"omitempty,phone"`
	Sex      *string `json:"sex" binding:"omitempty,sex"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var uri getUserUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidRequest))
		return
	}

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidRequest))
		return
	}

	if req.Fullname == nil && req.Email == nil && req.Phone == nil && req.Sex == nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("at least one field must be provided")))
		return
	}

	// user can update only his own data!!
	payloadVal, exists := ctx.Get(authPayloadKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrMissingAuthPayload))
		return
	}
	payload, ok := payloadVal.(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrMissingAuthPayload))
		return
	}
	if payload.Username != uri.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrPermissionDenied))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, uri.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServerError))
		return
	}

	args := db.UpdateUserParams{
		ID:       user.ID,
		Phone:    user.Phone,
		FullName: user.FullName,
		Sex:      user.Sex,
		Email:    user.Email,
	}
	if req.Fullname != nil {
		args.FullName = *req.Fullname
	}
	if req.Email != nil {
		args.Email = *req.Email
	}
	if req.Phone != nil {
		args.Phone = *req.Phone
	}
	if req.Sex != nil {
		args.Sex = *req.Sex
	}

	updated, err := server.store.UpdateUser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, newUserResponse(updated))
}
