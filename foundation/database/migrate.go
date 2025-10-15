package database

// Migrate runs database migrations
func (d *Database) Migrate() error {
	// Since we're using Docker init.sql, all tables are created automatically
	// No migrations needed - database schema is handled by Docker
	return nil
}
