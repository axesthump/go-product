package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func (r *CheckoutRepository) ClearCart(ctx context.Context, user int64) error {
	err := r.inTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		query := `SELECT id FROM carts WHERE user_id = $1`
		var cartID int64
		err := tx.QueryRow(ctx, query, user).Scan(&cartID)
		if err != nil {
			return err
		}

		err = r.deleteCart(ctx, tx, err, cartID)
		if err != nil {
			return err
		}

		err = r.deleteCartItems(ctx, tx, err, cartID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("postgres clear cart: %w", err)
	}
	return nil
}

func (r *CheckoutRepository) deleteCartItems(ctx context.Context, tx pgx.Tx, err error, cartID int64) error {
	queryDeleteCartItems := `DELETE FROM cart_items where cart_id = $1`
	_, err = tx.Exec(ctx, queryDeleteCartItems, cartID)
	if err != nil {
		return fmt.Errorf("postgres deleteCartItems: %w", err)
	}
	return nil
}

func (r *CheckoutRepository) deleteCart(ctx context.Context, tx pgx.Tx, err error, cartID int64) error {
	queryDeleteCart := `DELETE FROM carts where id = $1`

	_, err = tx.Exec(ctx, queryDeleteCart, cartID)
	if err != nil {
		return fmt.Errorf("postgres deleteCart: %w", err)
	}
	return nil
}
