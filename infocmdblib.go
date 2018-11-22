package infocmdbGoLib

type Config struct {
	Path string
	URL  string
}

type InfoCmdbGoLib struct {
	WS Webservice
	WC CmdbWebClient
}

func NewInfoCmdbGoLib() InfoCmdbGoLib {
	i := InfoCmdbGoLib{}

	return i
}

func (i *InfoCmdbGoLib) Login(url string, username string, password string) error {
	i.WC = NewCmdbWebClient()
	i.WS.client = &i.WC

	return i.WC.Login(url, username, password)
}
