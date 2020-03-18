package testing

const MetadataResponse = `
{
	"some_key": "some_value",
	"some_other_key": "some_other_value"
}
`

const SignalRequest = `
{
	"some_key": "some_value",
	"some_other_key": "some_other_value"
}
`

var (
	Metadata = map[string]interface{}{
		"some_key":       "some_value",
		"some_other_key": "some_other_value",
	}
)
