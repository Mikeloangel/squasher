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
func NewDbStorage(db *sql.DB) Storager {
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

	stm := `insert into links (short_id, original) values ($1,$2) returning id`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	row := ds.db.QueryRowContext(ctx, stm, short, url)
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

	stm := `select id, short_id, original from links where short_id=$1`

	row := ds.db.QueryRowContext(ctx, stm, short)
	err = row.Scan(&si.ID, &si.Shorten, &si.URL)
	return si, err
}

// Init opens file to write
func (s *DBStorage) Init() error {
	err := s.createSchemasIfNotExist()
	if err != nil {
		return err
	}

	return nil
}

// createSchemasIfNotExist creates links table
func (s *DBStorage) createSchemasIfNotExist() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stm := `
	create table if not exists links(
		id serial primary key,
		short_id varchar(50) not null UNIQUE,
		original VARCHAR(9000) not null UNIQUE
	);
	`

	_, err := s.db.ExecContext(ctx, stm)
	if err != nil {
		return err
	}
	return nil
}
