package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/Mikeloangel/squasher/internal/apperrors"
	urlgenerator "github.com/Mikeloangel/squasher/internal/urlGenerator"
)

// Storage represents db storage for shortened URLs.
type PostgresStorage struct {
	db        *sql.DB
	dbTimeout time.Duration
}

// NewPostgresStorage creates a new instance of file storage
func NewPostgresStorage(db *sql.DB, dbTimeout time.Duration) Storager {
	return &PostgresStorage{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

// Set stores the given URL and returns a shortened version of it.
// If the URL is already stored, it returns the existing shortened version.
func (ds *PostgresStorage) StoreURL(url string) (StorageItem, error) {
	short := urlgenerator.HashURL(url)

	si, err := ds.FetchURL(short)
	if err == nil {
		return si, apperrors.ErrItemAlreadyExists
	}

	si = StorageItem{
		URL:     url,
		Shorten: short,
	}

	query := `insert into links (short_id, original) values ($1,$2) returning id`
	ctx, cancel := context.WithTimeout(context.Background(), ds.dbTimeout)
	defer cancel()

	row := ds.db.QueryRowContext(ctx, query, short, url)
	err = row.Scan(&si.ID)
	if err != nil {
		return si, err
	}

	return si, nil
}

// Get retrieves the original URL for the given shortened version.
// It returns an error if the shortened URL is not found.
func (ds *PostgresStorage) FetchURL(short string) (StorageItem, error) {
	var err error
	si := StorageItem{}

	ctx, cancel := context.WithTimeout(context.Background(), ds.dbTimeout)
	defer cancel()

	query := `select id, short_id, original from links where short_id=$1`

	row := ds.db.QueryRowContext(ctx, query, short)
	err = row.Scan(&si.ID, &si.Shorten, &si.URL)
	return si, err
}

// Init opens file to write
func (ds *PostgresStorage) Init() error {
	err := ds.createSchemasIfNotExist()
	if err != nil {
		return err
	}

	return nil
}

// createSchemasIfNotExist creates links table
func (ds *PostgresStorage) createSchemasIfNotExist() error {
	ctx, cancel := context.WithTimeout(context.Background(), ds.dbTimeout)
	defer cancel()

	query := `
	create table if not exists links(
		id serial primary key,
		short_id varchar(50) not null,
		original VARCHAR(9000) not null,
		created_at timestamp not null default CURRENT_TIMESTAMP,
		UNIQUE (short_id, original)
	);
	`

	_, err := ds.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

// MultiStoreURL creates for slice of StorageItemOptionsInterface links and updates items
func (ds *PostgresStorage) MultiStoreURL(items *[]StorageItemOptionsInterface) error {
	ctx, cancel := context.WithTimeout(context.Background(), ds.dbTimeout)
	defer cancel()

	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = ds.storeMultiItemsInTransaction(ctx, tx, items)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (ds *PostgresStorage) storeMultiItemsInTransaction(ctx context.Context, tx *sql.Tx, items *[]StorageItemOptionsInterface) error {
	insertStatement, selectStatement, err := ds.prepareStatementsForMultiStore(ctx, tx)
	defer selectStatement.Close()
	defer insertStatement.Close()

	if err != nil {
		return err
	}

	for _, v := range *items {
		si := v.GetStorageItem()
		shorten := urlgenerator.HashURL(si.URL)
		row := insertStatement.QueryRowContext(ctx, shorten, si.URL)
		err := row.Scan(&si.ID, &si.Shorten)

		// on conflict select
		if err == sql.ErrNoRows {
			row = selectStatement.QueryRowContext(ctx, si.URL)
			err = row.Scan(&si.ID, &si.Shorten)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		si.Shorten = shorten
	}
	return nil
}

func (ds *PostgresStorage) prepareStatementsForMultiStore(ctx context.Context, tx *sql.Tx) (insertStatement, selectStatement *sql.Stmt, err error) {
	insertQuery := `
	insert into links (short_id, original)
		values($1, $2)
		ON CONFLICT (short_id, original)
		DO NOTHING
		returning id, short_id
	`
	selectQuery := `select id, short_id from links where original = $1`

	insertStatement, err = tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return nil, nil, err
	}

	selectStatement, err = tx.PrepareContext(ctx, selectQuery)
	if err != nil {
		insertStatement.Close()
		return nil, nil, err
	}

	return
}
