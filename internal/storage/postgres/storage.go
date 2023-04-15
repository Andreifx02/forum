package postrgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Andreifx02/forum/internal/config"
	"github.com/Andreifx02/forum/internal/domain"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostrgeSQL.Host, cfg.PostrgeSQL.Port, cfg.PostrgeSQL.User, cfg.PostrgeSQL.Password, cfg.PostrgeSQL.DbName,
	)

	pool, err := pgxpool.Connect(context.Background(), psqlConn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Storage{
		pool: pool,
	}, nil

}

func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := s.GetUserByName(ctx, user.Nickname)
	if err == nil {
		return UserAlreadyExists
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO users (id, nickname) VALUES ($1, $2)
	`, user.ID, user.Nickname)

	return err
}

func (s *Storage) GetUserByName(ctx context.Context, username string) (*domain.User, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT * FROM users WHERE nickname = $1
	`, username)
	var user domain.User
	err := row.Scan(&user.ID, &user.Nickname)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) CreatePost(ctx context.Context, post *domain.Post) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO posts (id, author_id, topic, text, date) VALUES ($1, $2, $3, $4, $5)
	`, post.ID, post.AuthorID, post.Topic, post.Text, post.Date)

	return err
}

func (s *Storage) CreateSubscription(ctx context.Context, subscription *domain.Subscriptions) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO subscriptions (user_id, sub_user_id) VALUES ($1, $2) 
	`, subscription.ID, subscription.SubID)

	return err
}

func (s *Storage) CreateLike(ctx context.Context, like *domain.Like) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO likes (user_id, post_id) VALUES ($1, $2)
	`, like.UserID, like.PostID)

	return err
}

func (s *Storage) GetSubFeed(ctx context.Context, userID uuid.UUID) ([]domain.Post, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT * FROM posts 
		WHERE author_id IN (
			SELECT sub_user_id FROM subscriptions 
			WHERE user_id = $1
		) 
		ORDER BY date DESC 
	`, userID)

	if err != nil {
		return nil, err
	}

	return fetchPosts(rows)
}

func (s *Storage) GetPossibleFriends(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT sub_user_id FROM subscriptions 
		WHERE 
			user_id IN (SELECT sub_user_id FROM subscriptions WHERE user_id = $1)
	`, userID)

	if err != nil {
		return nil, err
	}

	users_ID := make([]uuid.UUID, 0)
	for rows.Next() {
		var user_id uuid.UUID
		err = rows.Scan(&user_id)
		if err != nil {
			return nil, fmt.Errorf("Scan error: %w", err)
		}
		users_ID = append(users_ID, user_id)
	}
	return users_ID, nil
}

func (s *Storage) GetInteresting(ctx context.Context, userID uuid.UUID) ([]domain.Post, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT * FROM posts WHERE author_id IN (
			SELECT sub_user_id FROM subscriptions 
			WHERE 
				(user_id IN (SELECT sub_user_id FROM subscriptions WHERE user_id = $1)
			OR
			 	user_id = $1)
		)
	`, userID)

	if err != nil {
		return nil, err
	}

	return fetchPosts(rows)
}

func fetchPosts(rows pgx.Rows) ([]domain.Post, error) {
	posts := make([]domain.Post, 0)

	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Topic, &post.Text, &post.Date)
		if err != nil {
			return nil, fmt.Errorf("Scan error: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Storage) GetFilterSubFeed(ctx context.Context, userID uuid.UUID, filter *domain.Filters) ([]domain.Post, error) {
	var filters []any

	// filters = append(filters, filter.Topic)
	// filters = append(filters, filter.KeyWords)

	query := "SELECT * FROM posts WHERE"

	if filter.DateFrom != nil {
		filters = append(filters, *filter.DateFrom)
		query += " date > $" + fmt.Sprint(len(filters))
	}
	if filter.DateTo != nil {
		if len(filters) != 0 {
			query += " AND"
		}
		filters = append(filters, *filter.DateTo)
		query += " date < $" + fmt.Sprint(len(filters))
	}
	if filter.AuthorID != nil {
		if len(filters) != 0 {
			query += " AND"
		}
		filters = append(filters, *filter.AuthorID)
		query += " author_id = $" + fmt.Sprint(len(filters))
	}

	query += " ORDER BY date DESC"
	fmt.Print(query)
	rows, err := s.pool.Query(ctx, query, filters...)

	if err != nil {
		return nil, err
	}

	return fetchPosts(rows)
}
