package user

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/huydnt1801/chuyende/internal/ent"
	"github.com/huydnt1801/chuyende/internal/ent/user"
	"github.com/huydnt1801/chuyende/pkg/log"
)

type Repo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, u *User, updated *UserUpdate) error
	FindUser(ctx context.Context, params *UserParams) (*User, error)
}

type RepoImpl struct {
	logger logr.Logger
	client *ent.Client
}

func NewRepo(client *ent.Client) *RepoImpl {
	return &RepoImpl{
		logger: log.ZapLogger(),
		client: client,
	}
}

func (r *RepoImpl) CreateUser(ctx context.Context, user *User) (*User, error) {
	q := r.client.User.Create().
		SetPhoneNumber(user.PhoneNumber).
		SetPassword(user.Password).
		SetConfirmed(user.Confirmed)
	if s := user.FullName; s != "" {
		q.SetFullName(s)
	}
	u, err := q.Save(ctx)
	if ent.IsConstraintError(err) {
		return nil, UserExistError{}
	}
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	user.ID = u.ID
	return user, nil
}

func (r *RepoImpl) UpdateUser(ctx context.Context, u *User, updated *UserUpdate) (*User, error) {
	q := r.client.User.UpdateOneID(u.ID)
	if v := updated.FullName; v != "" {
		q.SetFullName(v)
	}
	if v := updated.Password; v != "" {
		q.SetPassword(v)
	}
	if v := updated.ImageURL; v != "" {
		q.SetImageURL(v)
	}
	if v := updated.Confirmed; v != nil {
		q.SetConfirmed(*v)
	}
	user, err := q.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed updating user: %w", err)
	}
	return DecodeUser(user)
}

func (r *RepoImpl) FindUser(ctx context.Context, params *UserParams) (*User, error) {
	q := r.client.User.Query()
	if v := params.PhoneNumber; v != "" {
		q = q.Where(user.PhoneNumber(v))
	}
	if v := params.ID; v != 0 {
		q = q.Where(user.ID(v))
	}

	u, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, UserNotFoundError{}
		}
		return nil, err
	}
	return DecodeUser(u)
}
