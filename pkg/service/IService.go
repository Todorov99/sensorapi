package service

type IService interface {
	GetAll() (interface{}, error)
	GetById(ID string) (interface{}, error)
	Add(model interface{}) error
	Update(model interface{}) error
	Delete(ID string) (interface{}, error)
}
