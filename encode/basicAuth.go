package encode

import "encoding/base64"

func BasicAuth(username, password string) string {
    auth := username + ":" + password
    return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
