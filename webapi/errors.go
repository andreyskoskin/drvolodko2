package webapi

type (
	Message struct {
		Message string `json:"message"`
	}

	NotFound struct {
		Entity string `json:"entity"`
		ID     int64  `json:"id"`
	}
)

var (
	updated = Message{
		Message: "Updated",
	}

	internalError = Message{
		Message: "Internal Server Error",
	}
)
