package hackbrowserdata

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/moond4rk/hackbrowserdata/browserdata"
	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
)

type chromium struct {
	name           browser
	storage        string
	profilePath    string
	enableAllUsers bool
	profilePaths   []string
	masterKey      []byte
	// defaultDataTypes
	supportedDataTypes []DataType
	extractors         map[DataType]browserdata.Extractor
	extractedData      map[DataType]interface{}
}

func NewChromium(options *Options) (Browser, error) {
	if options.ProfilePath == "" {
		return nil, errors.New("profile path is required")
	}
	if options.Name == "" {
		return nil, errors.New("browser name is required")
	}
	c := &chromium{
		name:               options.Name,
		profilePath:        options.ProfilePath,
		enableAllUsers:     true,
		supportedDataTypes: defaultDataTypes,
		extractors:         make(map[DataType]browserdata.Extractor),
		extractedData:      make(map[DataType]interface{}),
	}
	if !options.IsEnableAllUser {
		c.enableAllUsers = false
	}
	if len(options.DataTypes) > 0 {
		c.supportedDataTypes = options.DataTypes
	}
	if err := c.init(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *chromium) init() error {
	if err := c.initProfiles(); err != nil {
		return fmt.Errorf("profile path '%s' does not exist %w", c.profilePath, ErrBrowserNotExists)
	}
	if err := c.initExtractors(); err != nil {
		return err
	}
	return c.initMasterKey()
}

func (c *chromium) ExtractBrowserData(dataTypes []DataType) (map[DataType]interface{}, error) {
	for _, dataType := range dataTypes {
		if extractor, ok := c.extractors[dataType]; ok {
			data, err := extractor.Extract()
			if err != nil {
				fmt.Printf("extract %s data failed: %v", dataType, err)
				continue
			}
			c.extractedData[dataType] = data
		}
	}
	return c.extractedData, nil
}

// func (c *chromium) Passwords() ([]password.Password, error) {
// 	// browserData, err := c.ExtractBrowserData([]DataType{TypePassword})
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	dataType := TypePassword
// 	if data, ok := c.extractedData[dataType]; ok {
// 		return data.([]password.Password), nil
// 	}
// 	extractor, ok := c.extractors[dataType]
// 	if !ok {
// 		return nil, fmt.Errorf("%s extractor for %s not found", dataType, c.name)
// 	}
// 	data, err := extractor.ExtractChromium()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data.([]password.Password), nil
// }

func (c *chromium) filterExistDataPaths(dataTypes []DataType) (map[DataType][]string, error) {
	// exporters := make(map[DataType]BrowserData)
	dataPaths := make(map[DataType][]string)
	var errs []error
	for _, profile := range c.profilePaths {
		for _, dataType := range dataTypes {
			dataTypeFile := filepath.Join(profile, dataType.Filename(c.name))
			if !fileutil.IsFileExists(dataTypeFile) {
				errs = append(errs, ErrBrowsingDataNotExists)
			}
			dataPaths[dataType] = append(dataPaths[dataType], dataTypeFile)
		}
	}
	return dataPaths, nil
}

func (c *chromium) Passwords() ([]browserdata.Password, error) {
	dataType := TypePassword
	if data, ok := c.extractedData[dataType]; ok {
		return data.([]browserdata.Password), nil
	}
	extractor, ok := c.extractors[dataType]
	if !ok {
		return nil, fmt.Errorf("%s extractor for %s not found", dataType, c.name)
	}
	data, err := extractor.Extract()
	if err != nil {
		return nil, err
	}
	return data.([]browserdata.Password), nil
}

func (c *chromium) Cookies() ([]browserdata.Cookie, error) {
	dataType := TypeCookie
	if data, ok := c.extractedData[dataType]; ok {
		return data.([]browserdata.Cookie), nil
	}
	extractor, ok := c.extractors[dataType]
	if !ok {
		return nil, fmt.Errorf("%s extractor for %s not found", dataType, c.name)
	}
	data, err := extractor.Extract()
	if err != nil {
		return nil, err
	}
	return data.([]browserdata.Cookie), nil
}

func (c *chromium) initExtractors() error {
	dataPaths, err := c.filterExistDataPaths(c.supportedDataTypes)
	if err != nil {
		return err
	}
	for _, dataType := range c.supportedDataTypes {
		if _, ok := dataPaths[dataType]; !ok {
			continue
		}
		c.extractors[dataType] = dataType.NewExtractor(c.name.Type(), c.masterKey, dataPaths[dataType])
	}
	return nil
}

func (c *chromium) initProfiles() error {
	if !fileutil.IsDirExists(c.profilePath) {
		return ErrBrowserNotExists
	}
	if c.enableAllUsers {
		profilesPaths, err := c.findAllProfiles()
		if err != nil {
			return err
		}
		c.profilePaths = profilesPaths
	} else {
		c.profilePaths = []string{c.profilePath}
	}
	return nil
}

// TODO: mix it as firefox's find All Profiles
func (c *chromium) findAllProfiles() ([]string, error) {
	var profiles []string
	root := fileutil.ParentDir(c.profilePath)
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// if the path ends with "History", add it to the list
		if strings.HasSuffix(path, TypeHistory.Filename(c.name)) {
			// skip the "System Profile" directory
			if !strings.Contains(path, "System Profile") {
				profiles = append(profiles, filepath.Dir(path))
			}
		}

		// calculate the depth of the current path
		depth := len(strings.Split(path, string(filepath.Separator))) - len(strings.Split(root, string(filepath.Separator)))
		// if the depth is more than 2 and it's a directory, skip it
		if info.IsDir() && path != root && depth >= 2 {
			return filepath.SkipDir
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return profiles, err
}
