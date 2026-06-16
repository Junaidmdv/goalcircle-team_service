package postgres



func (db *postgresDB) Migration() error {
	if err := db.DB.AutoMigrate(); err != nil {
		return err
	}
	return nil
}
