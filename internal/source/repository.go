package source

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate() error {
	return r.db.AutoMigrate(&Source{})
}

func (r *Repository) Create(src *Source) error {
	return r.db.Create(src).Error
}

func (r *Repository) FindOne(src *Source) error {
	return r.db.First(&src).Error

}

func (r *Repository) Delete(src *Source) error {
	return r.db.Delete(&src).Error

}
