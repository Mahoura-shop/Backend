package enum

type BucketType uint

const (
	ProductPic BucketType = iota + 1
)

func (bt BucketType) String() string {
	switch bt {
	case ProductPic:
		return "productPic"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		ProductPic,
	}
}
