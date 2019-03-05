package resourcebundle

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	duplicateBundle = errors.New("bundle with the same language and local code already exist")
	Languages       = []Language{
		{"Amharic", "am", "AM"},
		{"Arabic", "ar", "AR"},
		{"Basque", "eu", "EU"},
		{"Bengali", "bn", "BN"},
		{"English (UK)", "en", "GB"},
		{"Portuguese (Brazil)", "pt", "BR"},
		{"Bulgarian", "bg", "BG"},
		{"Catalan", "ca", "CA"},
		{"Cherokee", "chr", "CHR"},
		{"Croatian", "hr", "HR"},
		{"Czech", "cs", "CS"},
		{"Danish", "da", "DA"},
		{"Dutch", "nl", "NL"},
		{"English (US)", "en", "EN"},
		{"Estonian", "et", "ET"},
		{"Filipino", "fil", "FIL"},
		{"Finnish", "fi", "FI"},
		{"French", "fr", "FR"},
		{"German", "de", "DE"},
		{"Greek", "el", "EL"},
		{"Gujarati", "gu", "GU"},
		{"Hebrew", "iw", "IW"},
		{"Hindi", "hi", "HI"},
		{"Hungarian", "hu", "HU"},
		{"Icelandic", "is", "IS"},
		{"Indonesian", "id", "ID"},
		{"Italian", "it", "IT"},
		{"Japanese", "ja", "JA"},
		{"Kannada", "kn", "KN"},
		{"Korean", "ko", "KO"},
		{"Latvian", "lv", "LV"},
		{"Lithuanian", "lt", "LT"},
		{"Malay", "ms", "MS"},
		{"Malayalam", "ml", "ML"},
		{"Marathi", "mr", "MR"},
		{"Norwegian", "no", "NO"},
		{"Polish", "pl", "PL"},
		{"Portuguese (Portugal)", "pt", "PT"},
		{"Romanian", "ro", "RO"},
		{"Russian", "ru", "RU"},
		{"Serbian", "sr", "SR"},
		{"Chinese (PRC)", "zh", "CN"},
		{"Slovak", "sk", "SK"},
		{"Slovenian", "sl", "SL"},
		{"Spanish", "es", "ES"},
		{"Swahili", "sw", "SW"},
		{"Swedish", "sv", "SV"},
		{"Tamil", "ta", "TA"},
		{"Telugu", "te", "TE"},
		{"Thai", "th", "TH"},
		{"Chinese (Taiwan)", "zh", "TW"},
		{"Turkish", "tr", "TR"},
		{"Urdu", "ur", "UR"},
		{"Ukrainian", "uk", "UK"},
		{"Vietnamese", "vi", "VI"},
		{"Welsh", "cy", "CY"},
	}
)

type Language struct {
	Name         string
	LanguageCode string
	LocalCode    string
}

// Bundle represent one resource bundle, belong to one language and local code
type Bundle struct {
	LanguageCode string            `json:"languageCode"`
	LocalCode    string            `json:"localCode"`
	TextMap      map[string]string `json:"textMap"`
}

// toProperties export the bundle content into properties KV style content.
func (b *Bundle) toProperties() []byte {
	var buffer bytes.Buffer
	for k, v := range b.TextMap {
		buffer.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	return buffer.Bytes()
}

// BundleFromPropertiesFile Create a Bundle from a properties file.
func BundleFromPropertiesFile(languageCode, localCode, path string) (*Bundle, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	} else {
		return BundleFromPropertiesData(languageCode, localCode, data)
	}
}

// BundleFromPropertiesData Create a Bundle from a content of properties file
func BundleFromPropertiesData(languageCode, localCode string, propertiesData []byte) (*Bundle, error) {
	b := &Bundle{
		LanguageCode: languageCode,
		LocalCode:    localCode,
	}
	err := b.fromProperties(propertiesData)
	if err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

// fromProperties import the content of the bundle from a file content in KV properties format.
func (b *Bundle) fromProperties(data []byte) error {
	b.TextMap = make(map[string]string)
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			k := line[:strings.Index(line, "=")]
			v := line[strings.Index(line, "=")+1:]
			b.TextMap[k] = v
		}
	}
	return scanner.Err()
}

// ResourceBundle represent a set of internalization resource bundle.
// it has the default languange and local code to look for values using Get method.
type ResourceBundle struct {
	LanguageCode string    `json:"languageCode"`
	LocalCode    string    `json:"localCode"`
	Default      *Bundle   `json:"default"`
	Bundles      []*Bundle `json:"bundles"`
}

// Import a resource bundle zip data content into ResourceBundle struct.
func ZipImport(language, local string, zipData []byte) (*ResourceBundle, error) {
	reader := bytes.NewReader(zipData)
	zipReader, err := zip.NewReader(reader, int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	ret := &ResourceBundle{
		LanguageCode: language,
		LocalCode:    local,
		Bundles:      make([]*Bundle, 0),
	}

	var defLange string
	var defLocal string

	for _, f := range zipReader.File {
		name := f.Name
		if name == "meta.properties" {
			props, err := processPropertiesFile(f)
			if err != nil {
				return nil, err
			}
			defLange = props["defaultLang"]
			defLocal = props["defaultLocal"]
		} else {
			langPart := name[:strings.Index(name, ".")]
			strArr := strings.Split(langPart, "_")
			props, err := processPropertiesFile(f)
			if err != nil {
				return nil, err
			}
			b := &Bundle{
				LocalCode:    strArr[1],
				LanguageCode: strArr[0],
				TextMap:      props,
			}
			err = ret.AddBundle(b, false)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, b := range ret.Bundles {
		if b.LocalCode == defLocal && b.LanguageCode == defLange {
			ret.Default = b
		}
	}
	return ret, nil
}

// processPropertiesFile read a data in KV properties format into map of string to string.
func processPropertiesFile(f *zip.File) (map[string]string, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	ret := make(map[string]string)
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			k := line[:strings.Index(line, "=")]
			v := line[strings.Index(line, "=")+1:]
			ret[k] = v
		}
	}
	return ret, nil
}

// ToJsonString export this ResourceBundle into JSON
func (rb *ResourceBundle) ToJsonString() ([]byte, error) {
	return json.Marshal(rb)
}

// FromJsonString will create a new ResourceBundle from JSON
func FromJsonString(jsondata []byte) (*ResourceBundle, error) {
	rb := &ResourceBundle{}
	err := json.Unmarshal(jsondata, rb)
	if err != nil {
		return nil, err
	} else {
		defLang := rb.Default.LanguageCode
		defLoca := rb.Default.LocalCode
		for _, b := range rb.Bundles {
			if b.LocalCode == defLoca && b.LanguageCode == defLang {
				rb.Default = b
				break
			}
		}
		return rb, err
	}
}

// NewResourceBundle create new ResourceBundle instance
func NewResourceBundle(language, local string, defaultBundle *Bundle, bundles []*Bundle) *ResourceBundle {
	rb := &ResourceBundle{
		LanguageCode: language,
		LocalCode:    local,
		Default:      defaultBundle,
		Bundles:      bundles,
	}
	return rb
}

// ZipExport export the ResourceBundle into bytes in ZIP format.
func (rb *ResourceBundle) ZipExport() ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// make the meta file
	metabody := fmt.Sprintf("defaultLang=%s\ndefaultLocal=%s", rb.Default.LanguageCode, rb.Default.LocalCode)

	f, err := w.Create("meta.properties")
	if err != nil {
		return nil, err
	}
	_, err = f.Write([]byte(metabody))
	if err != nil {
		return nil, err
	}

	for _, b := range rb.Bundles {
		f, err := w.Create(fmt.Sprintf("%s_%s.properties", b.LanguageCode, b.LocalCode))
		if err != nil {
			return nil, err
		}
		_, err = f.Write(b.toProperties())
		if err != nil {
			return nil, err
		}
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// AddBundle add a language bundle into Resource Bundle
func (rb *ResourceBundle) AddBundle(bundle *Bundle, isDefault bool) error {
	for _, b := range rb.Bundles {
		if b.LocalCode == bundle.LocalCode && b.LanguageCode == bundle.LanguageCode {
			return duplicateBundle
		}
	}
	rb.Bundles = append(rb.Bundles, bundle)
	if isDefault {
		rb.Default = bundle
	}
	return nil
}

// Get a value from key in the resource bundle. First it will look for key on the languange and local as defined in the ResourceBundle struct.
// if the key is not found, it will use the dafault bundle to look for one. If still not found, it will return an empty string.
func (rb *ResourceBundle) Get(key string) string {
	bundle := rb.GetBundle(rb.LanguageCode, rb.LocalCode)
	if bundle == nil {
		bundle = rb.Default
		if bundle == nil {
			return ""
		}
	}
	if txt, ok := bundle.TextMap[key]; ok {
		return txt
	} else {
		if bundle == rb.Default {
			return ""
		}
		if txt, ok := rb.Default.TextMap[key]; ok {
			return txt
		} else {
			return ""
		}
	}
}

// GetBundle look for Bundle instance in ResourceBundle with specified language and local. If not found, it will look for bundle with same language regardless of its local.
func (rb *ResourceBundle) GetBundle(languageCode, localCode string) *Bundle {
	for _, b := range rb.Bundles {
		if b.LanguageCode == rb.LanguageCode && b.LocalCode == rb.LocalCode {
			return b
		}
	}
	for _, b := range rb.Bundles {
		if b.LanguageCode == rb.LanguageCode {
			return b
		}
	}
	return nil
}
