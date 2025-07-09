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

func (r *Repository) FindByID(id uint) (*Source, error) {
	var src Source
	if err := r.db.First(&src, id).Error; err != nil {
		return nil, err
	}
	return &src, nil
}

func (r *Repository) FindAll(sources *[]*Source) {
	r.db.Find(sources)

}

func (r *Repository) Delete(src *Source) error {
	return r.db.Delete(&src).Error

}
