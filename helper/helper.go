package helper

type saltDataType []byte

const (
	//Create set password byte length to 200
	pwByteLen = 200
)

/*
 *Generate random data for a given len
 *
 *len: the length of the data
 *
 *
 *return: byte of data
 */
func generateSalt() saltDataType {
	return saltDataType(make([]byte, 200))
}
