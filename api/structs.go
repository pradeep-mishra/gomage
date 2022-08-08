package api

type QueryParams struct {
	Flip      string  `query:"flip" validate:"omitempty,oneof=h v"`
	Pixlate   uint8   `query:"pixlate" validate:"omitempty,gte=1,lte=100"`
	SmartCrop string  `query:"smartcrop"`
	Sharpen   string  `query:"sharpen"`
	Rotate    string  `query:"rotate" validate:"omitempty,oneof=0 90 180 270"`
	Scale     float64 `query:"scale" validate:"omitempty,min=0,max=10"`
	Repeat    string  `query:"repeat"`
	Modulate  string  `query:"modulate"`
	Label     string  `query:"label"`
	Zoom      int     `query:"zoom" validate:"omitempty,min=1,max=10"`
	Format    string  `query:"format" validate:"omitempty,oneof=jpg jpeg png webp"`
}

type SmartCrop struct {
	Width    int    `validate:"omitempty,gt=1"`
	Height   int    `validate:"omitempty,gt=1"`
	CropType string `validate:"omitempty,oneof=entropy centre center attention low high all last"`
}

type Sharpen struct {
	Sigma     float64 `validate:"omitempty,gt=1"`
	Threshold float64 `validate:"omitempty,gt=1"`
	Slope     float64 `validate:"omitempty,gt=1"`
}
