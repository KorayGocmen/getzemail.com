package main

import "gorm.io/gorm"

func (i *MailInbox) AfterCreate(tx *gorm.DB) (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		var mail Mail
		if err := tx.First(&mail, "id = ?", i.MailID).Error; err != nil {
			return err
		}

		if err := tx.Model(&mail).Update("version", mail.Version+1).Error; err != nil {
			return err
		}

		return nil
	})
}

func (i *MailInbox) AfterDelete(tx *gorm.DB) (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		var mail Mail
		if err := tx.First(&mail, "id = ?", i.MailID).Error; err != nil {
			return err
		}

		if err := tx.Model(&mail).Update("version", mail.Version+1).Error; err != nil {
			return err
		}

		return nil
	})
}
