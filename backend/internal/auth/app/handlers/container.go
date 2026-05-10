package handlers

type Handlers struct {
	DummyLogin *DummyLoginHandler
}

func BuildHandlers(jwtToken string) *Handlers {
	return &Handlers{
		DummyLogin: NewDummyLoginHandler(jwtToken),
	}
}
