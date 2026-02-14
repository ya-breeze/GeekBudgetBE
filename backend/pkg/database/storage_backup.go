package database

func (s *storage) Backup(destination string) error {
	return s.db.Exec("VACUUM INTO ?", destination).Error
}
