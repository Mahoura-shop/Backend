package enum

type BucketType uint

const (
	ProductPic BucketType = iota + 1
	CategoryPic
	BrandPic
)

func (bt BucketType) String() string {
	switch bt {
	case ProductPic:
		return "productPic"
	case CategoryPic:
		return "categoryPic"
	case BrandPic:
		return "brandPic"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		ProductPic,
		CategoryPic,
		BrandPic,
	}
}
