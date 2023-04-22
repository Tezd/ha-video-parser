package translate

type (
	Request struct {
		Paths []string `json:"paths"`
	}
)

func BatchTranslate(files []string) map[string]string {
	return nil
}
