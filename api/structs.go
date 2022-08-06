package api

type QueryParams struct {
	Flip      string `query:"flip" validate:"omitempty,oneof=h v"`
	Pixlate   uint8  `query:"pixlate" validate:"omitempty,gte=1,lte=100"`
	SmartCrop string `query:"smartcrop"`
}

type SmartCrop struct {
	Width    int    `validate:"omitempty,gt=0"`
	Height   int    `validate:"omitempty,gt=0"`
	CropType string `validate:"omitempty,oneof=entropy centre center attention low high all last"`
}
