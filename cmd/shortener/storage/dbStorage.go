package storage

import (
	"context"
	"database/sql"

	urlgenerator "github.com/Mikeloangel/squasher/internal/urlGenerator"
)

// Storage represents db storage for shortened URLs.
type DBStorage struct {
	db *sql.DB
}

// NewDbStorage creates a new instance of file storage
func NewDBStorage(db *sql.DB) Storager {
	return &DBStorage{
		db: db,
	}
}

// Set stores the given URL and returns a shortened version of it.
// If the URL is already stored, it returns the existing shortened version.
func (ds *DBStorage) StoreURL(url string) (StorageItem, error) {
	short := urlgenerator.HashURL(url)

	si, err := ds.FetchURL(short)
	if err == nil {
		return si, nil
	}

	si = StorageItem{
		URL:     url,
		Shorten: short,
	}

	query := `insert into links (short_id, original) values ($1,$2) returning id`
	ctx, cancel := context.WithCancel(context.Background())
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
func (ds *DBStorage) FetchURL(short string) (StorageItem, error) {
	var err error
	si := StorageItem{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `select id, short_id, original from links where short_id=$1`

	row := ds.db.QueryRowContext(ctx, query, short)
	err = row.Scan(&si.ID, &si.Shorten, &si.URL)
	return si, err
}

// Init opens file to write
func (ds *DBStorage) Init() error {
	err := ds.createSchemasIfNotExist()
	if err != nil {
		return err
	}

	return nil
}

// createSchemasIfNotExist creates links table
func (ds *DBStorage) createSchemasIfNotExist() error {
	ctx, cancel := context.WithCancel(context.Background())
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
func (ds *DBStorage) MultiStoreURL(items *[]StorageItemOptionsInterface) error {
	query := `
	insert into links (short_id, original) 
		values($1, $2)
		ON CONFLICT (short_id, original) 
    	DO NOTHING
		returning id, short_id
	`
	selectQuery := `select id, short_id from links where original = $1`

	tx, err := ds.db.Begin()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	selectStmt, err := tx.PrepareContext(ctx, selectQuery)
	if err != nil {
		return err
	}
	defer selectStmt.Close()

	for _, v := range *items {
		si := v.GetStorageItem()
		shorten := urlgenerator.HashURL(si.URL)
		row := stmt.QueryRowContext(ctx, shorten, si.URL)
		err := row.Scan(&si.ID, &si.Shorten)

		// on conflict select
		if err == sql.ErrNoRows {
			row = selectStmt.QueryRowContext(ctx, si.URL)
			err = row.Scan(&si.ID, &si.Shorten)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else if err != nil {
			tx.Rollback()
			return err
		}
		si.Shorten = shorten
	}

	tx.Commit()
	return nil
}
