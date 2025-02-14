package  repositoryImplementation


import (
	"errors"

	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/repository"
	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) repository.EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) GetByID(id string) (*model.Event, error){
	var event model.Event
	result := r.db.Where("id = ?", id).First(&event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound){
			return nil, nil
		}
		return nil, result.Error
	}
	return &event, nil
}

func (r *eventRepository) List(limit, offset int) ([]*model.Event, error) {
	var events []*model.Event
	result := r.db.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&events)
	
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

func (r *eventRepository) Create(event *model.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) Update(event *model.Event) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Save(event).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *eventRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&model.Event{}, "id = ?", id).Error
		if err != nil{
			return err
		}
		return nil
	})
}


