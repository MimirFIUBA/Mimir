package datahandler

type DataCollection []Data

type Data interface {
	GetId() string
	Update(updatedData Data)
}

func (elements DataCollection) GetElement(id string) *Data {
	for i := range elements {
		element := elements[i]
		if element.GetId() == id {
			return &elements[i]
		}
	}
	return nil
}
