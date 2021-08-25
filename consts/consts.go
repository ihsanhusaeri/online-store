package consts

type ResponseMessage string

const (
	InternalServerErrorMessage ResponseMessage = "Internal server error"
	BadRequestMessage                          = "Bad request"
	CreatedMessage                             = "Data created"
	SuccessMessage                             = "Success"
)

type OrderStatus string

const (
	Cart     OrderStatus = "cart"
	Checkout             = "checkout"
	Paid                 = "paid"
	Expired              = "expired"
)
