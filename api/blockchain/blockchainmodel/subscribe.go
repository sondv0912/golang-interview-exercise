package blockchainmodel

type SubscribeRequestBody struct {
	Address string `json:"address" bson:"address" validate:"required"`
}
