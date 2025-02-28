package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
)

// Options hold the global configuration values without further validation or processing.
// Application code should retrieve option values via getter functions since they provide
// validation and return defaults if a value is empty.
type Options struct {
	Name                  string        `json:"-"`
	Edition               string        `json:"-"`
	Version               string        `json:"-"`
	Copyright             string        `json:"-"`
	PartnerID             string        `yaml:"-" json:"-" flag:"partner-id"`
	AuthMode              string        `yaml:"AuthMode" json:"-" flag:"auth-mode"`
	Public                bool          `yaml:"Public" json:"-" flag:"public"`
	AdminUser             string        `yaml:"AdminUser" json:"-" flag:"admin-user"`
	AdminPassword         string        `yaml:"AdminPassword" json:"-" flag:"admin-password"`
	SessMaxAge            int64         `yaml:"SessMaxAge" json:"-" flag:"sess-maxage"`
	SessTimeout           int64         `yaml:"SessTimeout" json:"-" flag:"sess-timeout"`
	LogLevel              string        `yaml:"LogLevel" json:"-" flag:"log-level"`
	Prod                  bool          `yaml:"Prod" json:"Prod" flag:"prod"`
	Debug                 bool          `yaml:"Debug" json:"Debug" flag:"debug"`
	Trace                 bool          `yaml:"Trace" json:"Trace" flag:"trace"`
	Test                  bool          `yaml:"-" json:"Test,omitempty" flag:"test"`
	Unsafe                bool          `yaml:"-" json:"-" flag:"unsafe"`
	Demo                  bool          `yaml:"Demo" json:"-" flag:"demo"`
	Sponsor               bool          `yaml:"-" json:"-" flag:"sponsor"`
	ReadOnly              bool          `yaml:"ReadOnly" json:"ReadOnly" flag:"read-only"`
	Experimental          bool          `yaml:"Experimental" json:"Experimental" flag:"experimental"`
	ConfigPath            string        `yaml:"ConfigPath" json:"-" flag:"config-path"`
	DefaultsYaml          string        `json:"-" yaml:"-" flag:"defaults-yaml"`
	OriginalsPath         string        `yaml:"OriginalsPath" json:"-" flag:"originals-path"`
	OriginalsLimit        int           `yaml:"OriginalsLimit" json:"OriginalsLimit" flag:"originals-limit"`
	ResolutionLimit       int           `yaml:"ResolutionLimit" json:"ResolutionLimit" flag:"resolution-limit"`
	StoragePath           string        `yaml:"StoragePath" json:"-" flag:"storage-path"`
	SidecarPath           string        `yaml:"SidecarPath" json:"-" flag:"sidecar-path"`
	BackupPath            string        `yaml:"BackupPath" json:"-" flag:"backup-path"`
	CachePath             string        `yaml:"CachePath" json:"-" flag:"cache-path"`
	ImportPath            string        `yaml:"ImportPath" json:"-" flag:"import-path"`
	ImportDest            string        `yaml:"ImportDest" json:"-" flag:"import-dest"`
	AssetsPath            string        `yaml:"AssetsPath" json:"-" flag:"assets-path"`
	CustomAssetsPath      string        `yaml:"-" json:"-" flag:"custom-assets-path"`
	TempPath              string        `yaml:"TempPath" json:"-" flag:"temp-path"`
	Workers               int           `yaml:"Workers" json:"Workers" flag:"workers"`
	WakeupInterval        time.Duration `yaml:"WakeupInterval" json:"WakeupInterval" flag:"wakeup-interval"`
	AutoIndex             int           `yaml:"AutoIndex" json:"AutoIndex" flag:"auto-index"`
	AutoImport            int           `yaml:"AutoImport" json:"AutoImport" flag:"auto-import"`
	DisableWebDAV         bool          `yaml:"DisableWebDAV" json:"DisableWebDAV" flag:"disable-webdav"`
	DisableBackups        bool          `yaml:"DisableBackups" json:"DisableBackups" flag:"disable-backups"`
	DisableSettings       bool          `yaml:"DisableSettings" json:"-" flag:"disable-settings"`
	DisablePlaces         bool          `yaml:"DisablePlaces" json:"DisablePlaces" flag:"disable-places"`
	DisableTensorFlow     bool          `yaml:"DisableTensorFlow" json:"DisableTensorFlow" flag:"disable-tensorflow"`
	DisableFaces          bool          `yaml:"DisableFaces" json:"DisableFaces" flag:"disable-faces"`
	DisableClassification bool          `yaml:"DisableClassification" json:"DisableClassification" flag:"disable-classification"`
	DisableFFmpeg         bool          `yaml:"DisableFFmpeg" json:"DisableFFmpeg" flag:"disable-ffmpeg"`
	DisableExifTool       bool          `yaml:"DisableExifTool" json:"DisableExifTool" flag:"disable-exiftool"`
	DisableHeifConvert    bool          `yaml:"DisableHeifConvert" json:"DisableHeifConvert" flag:"disable-heifconvert"`
	DisableDarktable      bool          `yaml:"DisableDarktable" json:"DisableDarktable" flag:"disable-darktable"`
	DisableRawtherapee    bool          `yaml:"DisableRawtherapee" json:"DisableRawtherapee" flag:"disable-rawtherapee"`
	DisableSips           bool          `yaml:"DisableSips" json:"DisableSips" flag:"disable-sips"`
	DisableRaw            bool          `yaml:"DisableRaw" json:"DisableRaw" flag:"disable-raw"`
	RawPresets            bool          `yaml:"RawPresets" json:"RawPresets" flag:"raw-presets"`
	ExifBruteForce        bool          `yaml:"ExifBruteForce" json:"ExifBruteForce" flag:"exif-bruteforce"`
	DetectNSFW            bool          `yaml:"DetectNSFW" json:"DetectNSFW" flag:"detect-nsfw"`
	UploadNSFW            bool          `yaml:"UploadNSFW" json:"-" flag:"upload-nsfw"`
	DefaultTheme          string        `yaml:"DefaultTheme" json:"DefaultTheme" flag:"default-theme"`
	DefaultLocale         string        `yaml:"DefaultLocale" json:"DefaultLocale" flag:"default-locale"`
	AppIcon               string        `yaml:"AppIcon" json:"AppIcon" flag:"app-icon"`
	AppName               string        `yaml:"AppName" json:"AppName" flag:"app-name"`
	AppMode               string        `yaml:"AppMode" json:"AppMode" flag:"app-mode"`
	Imprint               string        `yaml:"Imprint" json:"Imprint" flag:"imprint"`
	ImprintUrl            string        `yaml:"ImprintUrl" json:"ImprintUrl" flag:"imprint-url"`
	WallpaperUri          string        `yaml:"WallpaperUri" json:"WallpaperUri" flag:"wallpaper-uri"`
	CdnUrl                string        `yaml:"CdnUrl" json:"CdnUrl" flag:"cdn-url"`
	SiteUrl               string        `yaml:"SiteUrl" json:"SiteUrl" flag:"site-url"`
	SiteAuthor            string        `yaml:"SiteAuthor" json:"SiteAuthor" flag:"site-author"`
	SiteTitle             string        `yaml:"SiteTitle" json:"SiteTitle" flag:"site-title"`
	SiteCaption           string        `yaml:"SiteCaption" json:"SiteCaption" flag:"site-caption"`
	SiteDescription       string        `yaml:"SiteDescription" json:"SiteDescription" flag:"site-description"`
	SitePreview           string        `yaml:"SitePreview" json:"SitePreview" flag:"site-preview"`
	TrustedProxies        []string      `yaml:"TrustedProxies" json:"-" flag:"trusted-proxy"`
	HttpMode              string        `yaml:"HttpMode" json:"-" flag:"http-mode"`
	HttpCompression       string        `yaml:"HttpCompression" json:"-" flag:"http-compression"`
	HttpHost              string        `yaml:"HttpHost" json:"-" flag:"http-host"`
	HttpPort              int           `yaml:"HttpPort" json:"-" flag:"http-port"`
	AutoTLS               string        `yaml:"AutoTLS" json:"AutoTLS" flag:"auto-tls"` // AutoTLS enabled automatic HTTPS via Let's Encrypt if set a valid email address.
	TLSKey                string        `yaml:"TLSKey" json:"TLSKey" flag:"tls-key"`
	TLSCert               string        `yaml:"TLSCert" json:"TLSCert" flag:"tls-cert"`
	HttpsPort             int           `yaml:"HttpsPort" json:"HttpsPort" flag:"https-port"` // HttpsPort is the port number to be used for HTTPS connections.
	HttpsProxyHeaders     []string      `yaml:"HttpsProxyHeaders" json:"-" flag:"https-proxy-header"`
	HttpsProxyProto       []string      `yaml:"HttpsProxyProto" json:"-" flag:"https-proxy-proto"`
	HttpsRedirect         int           `yaml:"HttpsRedirect" json:"HttpsRedirect" flag:"https-redirect"`
	DatabaseDriver        string        `yaml:"DatabaseDriver" json:"-" flag:"database-driver"`
	DatabaseDsn           string        `yaml:"DatabaseDsn" json:"-" flag:"database-dsn"`
	DatabaseName          string        `yaml:"DatabaseName" json:"-" flag:"database-name"`
	DatabaseServer        string        `yaml:"DatabaseServer" json:"-" flag:"database-server"`
	DatabaseUser          string        `yaml:"DatabaseUser" json:"-" flag:"database-user"`
	DatabasePassword      string        `yaml:"DatabasePassword" json:"-" flag:"database-password"`
	DatabaseConns         int           `yaml:"DatabaseConns" json:"-" flag:"database-conns"`
	DatabaseConnsIdle     int           `yaml:"DatabaseConnsIdle" json:"-" flag:"database-conns-idle"`
	DarktableBin          string        `yaml:"DarktableBin" json:"-" flag:"darktable-bin"`
	DarktableCachePath    string        `yaml:"DarktableCachePath" json:"-" flag:"darktable-cache-path"`
	DarktableConfigPath   string        `yaml:"DarktableConfigPath" json:"-" flag:"darktable-config-path"`
	DarktableBlacklist    string        `yaml:"DarktableBlacklist" json:"-" flag:"darktable-blacklist"`
	RawtherapeeBin        string        `yaml:"RawtherapeeBin" json:"-" flag:"rawtherapee-bin"`
	RawtherapeeBlacklist  string        `yaml:"RawtherapeeBlacklist" json:"-" flag:"rawtherapee-blacklist"`
	SipsBin               string        `yaml:"SipsBin" json:"-" flag:"sips-bin"`
	SipsBlacklist         string        `yaml:"SipsBlacklist" json:"-" flag:"sips-blacklist"`
	HeifConvertBin        string        `yaml:"HeifConvertBin" json:"-" flag:"heifconvert-bin"`
	FFmpegBin             string        `yaml:"FFmpegBin" json:"-" flag:"ffmpeg-bin"`
	FFmpegEncoder         string        `yaml:"FFmpegEncoder" json:"FFmpegEncoder" flag:"ffmpeg-encoder"`
	FFmpegBitrate         int           `yaml:"FFmpegBitrate" json:"FFmpegBitrate" flag:"ffmpeg-bitrate"`
	ExifToolBin           string        `yaml:"ExifToolBin" json:"-" flag:"exiftool-bin"`
	DetachServer          bool          `yaml:"DetachServer" json:"-" flag:"detach-server"`
	DownloadToken         string        `yaml:"DownloadToken" json:"-" flag:"download-token"`
	PreviewToken          string        `yaml:"PreviewToken" json:"-" flag:"preview-token"`
	ThumbColor            string        `yaml:"ThumbColor" json:"ThumbColor" flag:"thumb-color"`
	ThumbFilter           string        `yaml:"ThumbFilter" json:"ThumbFilter" flag:"thumb-filter"`
	ThumbSize             int           `yaml:"ThumbSize" json:"ThumbSize" flag:"thumb-size"`
	ThumbSizeUncached     int           `yaml:"ThumbSizeUncached" json:"ThumbSizeUncached" flag:"thumb-size-uncached"`
	ThumbUncached         bool          `yaml:"ThumbUncached" json:"ThumbUncached" flag:"thumb-uncached"`
	JpegQuality           string        `yaml:"JpegQuality" json:"JpegQuality" flag:"jpeg-quality"`
	JpegSize              int           `yaml:"JpegSize" json:"JpegSize" flag:"jpeg-size"`
	FaceSize              int           `yaml:"-" json:"-" flag:"face-size"`
	FaceScore             float64       `yaml:"-" json:"-" flag:"face-score"`
	FaceOverlap           int           `yaml:"-" json:"-" flag:"face-overlap"`
	FaceClusterSize       int           `yaml:"-" json:"-" flag:"face-cluster-size"`
	FaceClusterScore      int           `yaml:"-" json:"-" flag:"face-cluster-score"`
	FaceClusterCore       int           `yaml:"-" json:"-" flag:"face-cluster-core"`
	FaceClusterDist       float64       `yaml:"-" json:"-" flag:"face-cluster-dist"`
	FaceMatchDist         float64       `yaml:"-" json:"-" flag:"face-match-dist"`
	PIDFilename           string        `yaml:"PIDFilename" json:"-" flag:"pid-filename"`
	LogFilename           string        `yaml:"LogFilename" json:"-" flag:"log-filename"`
}

// NewOptions creates a new configuration entity by using two methods:
//
// 1. Load: This will initialize options from a yaml config file.
//
//  2. ApplyCliContext: Which comes after Load and overrides
//     any previous options giving an option two override file configs through the CLI.
func NewOptions(ctx *cli.Context) *Options {
	c := &Options{}

	// Has context?
	if ctx == nil {
		return c
	}

	// Set app name from metadata if possible.
	if s, ok := ctx.App.Metadata["Name"]; ok {
		c.Name = fmt.Sprintf("%s", s)
	}

	// Set app edition from metadata if possible.
	if s, ok := ctx.App.Metadata["Edition"]; ok {
		c.Edition = fmt.Sprintf("%s", s)
	}

	// Set copyright and version information.
	c.Copyright = ctx.App.Copyright
	c.Version = ctx.App.Version

	// Load defaults from YAML file?
	if defaultsYaml := ctx.GlobalString("defaults-yaml"); defaultsYaml == "" {
		log.Tracef("config: defaults yaml file not specified")
	} else if c.DefaultsYaml = fs.Abs(defaultsYaml); !fs.FileExists(c.DefaultsYaml) {
		log.Tracef("config: defaults file %s does not exist", clean.Log(c.DefaultsYaml))
	} else if err := c.Load(c.DefaultsYaml); err != nil {
		log.Warnf("config: failed loading defaults from %s (%s)", clean.Log(c.DefaultsYaml), err)
	}

	if err := c.ApplyCliContext(ctx); err != nil {
		log.Error(err)
	}

	return c
}

// expandFilenames converts path in config to absolute path
func (c *Options) expandFilenames() {
	c.ConfigPath = fs.Abs(c.ConfigPath)
	c.StoragePath = fs.Abs(c.StoragePath)
	c.BackupPath = fs.Abs(c.BackupPath)
	c.AssetsPath = fs.Abs(c.AssetsPath)
	c.CachePath = fs.Abs(c.CachePath)
	c.OriginalsPath = fs.Abs(c.OriginalsPath)
	c.ImportPath = fs.Abs(c.ImportPath)
	c.TempPath = fs.Abs(c.TempPath)
	c.PIDFilename = fs.Abs(c.PIDFilename)
	c.LogFilename = fs.Abs(c.LogFilename)
}

// Load uses a yaml config file to initiate the configuration entity.
func (c *Options) Load(fileName string) error {
	if fileName == "" {
		return nil
	}

	if !fs.FileExists(fileName) {
		return errors.New(fmt.Sprintf("%s not found", fileName))
	}

	yamlConfig, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlConfig, c)
}

// ApplyCliContext uses options from the CLI to setup configuration overrides
// for the entity.
func (c *Options) ApplyCliContext(ctx *cli.Context) error {
	return ApplyCliContext(c, ctx)
}
