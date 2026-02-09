package enum

type BucketType uint

const (
	ProductPic BucketType = iota + 1
	CategoryPic
)

func (bt BucketType) String() string {
	switch bt {
	case ProductPic:
		return "productPic"
	case CategoryPic:
		return "categoryPic"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		ProductPic,
		CategoryPic,
	}
}
