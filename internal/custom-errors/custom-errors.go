package customerrors

type OrderNotFound struct {
}

func (err OrderNotFound) Error() string {
	return "order not found"
}

type OrderAlreadyExist struct {
}

func (err OrderAlreadyExist) Error() string {
	return "order already exist"
}

type WrongID struct {
	Msg string
}

func (err WrongID) Error() string {
	return err.Msg
}