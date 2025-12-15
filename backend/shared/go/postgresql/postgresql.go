package postgresql

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPool(ctx context.Context, address string) (pool *pgxpool.Pool, err error) {
	for i := range 3 {
		pool, err = pgxpool.New(ctx, address)
		if err == nil {
			break
		}
		log.Printf("DB not ready yet (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return pool, err
}

//==============

// pool, err := postgresql.GetPool(context.Background(), os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		return fmt.Errorf("failed to connect db: %v", err)
// 	}
// 	defer pool.Close()
