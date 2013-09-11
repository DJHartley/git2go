package git

/*
#cgo pkg-config: libgit2
#include <git2.h>
#include <git2/errors.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type ConfigLevel int

const (
	// System-wide configuration file; /etc/gitconfig on Linux systems
	ConfigLevelSystem ConfigLevel = C.GIT_CONFIG_LEVEL_SYSTEM

	// XDG compatible configuration file; typically ~/.config/git/config
	ConfigLevelXDG ConfigLevel = C.GIT_CONFIG_LEVEL_XDG

	// User-specific configuration file (also called Global configuration
	// file); typically ~/.gitconfig
	ConfigLevelGlobal ConfigLevel = C.GIT_CONFIG_LEVEL_GLOBAL

	// Repository specific configuration file; $WORK_DIR/.git/config on
	// non-bare repos
	ConfigLevelLocal ConfigLevel = C.GIT_CONFIG_LEVEL_LOCAL

	// Application specific configuration file; freely defined by applications
	ConfigLevelApp ConfigLevel = C.GIT_CONFIG_LEVEL_APP

	// Represents the highest level available config file (i.e. the most
	// specific config file available that actually is loaded)
	ConfigLevelHighest ConfigLevel = C.GIT_CONFIG_HIGHEST_LEVEL
)


type Config struct {
	ptr *C.git_config
}

// NewConfig creates a new empty configuration object
func NewConfig() (*Config, error) {
	config := new(Config)

	ret := C.git_config_new(&config.ptr)
	if ret < 0 {
		return nil, LastError()
	}

	return config, nil
}

// AddFile adds a file-backed backend to the config object at the specified level.
func (c *Config) AddFile(path string, level ConfigLevel, force bool) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	ret := C.git_config_add_file_ondisk(c.ptr, cpath, C.git_config_level_t(level), cbool(force))
	if ret < 0 {
		return LastError()
	}

	return nil
}

func (c *Config) LookupInt32(name string) (int32, error) {
	var out C.int32_t
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.git_config_get_int32(&out, c.ptr, cname)
	if ret < 0 {
		return 0, LastError()
	}

	return int32(out), nil
}

func (c *Config) LookupInt64(name string) (int64, error) {
	var out C.int64_t
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.git_config_get_int64(&out, c.ptr, cname)
	if ret < 0 {
		return 0, LastError()
	}

	return int64(out), nil
}

func (c *Config) LookupString(name string) (string, error) {
	var ptr *C.char
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.git_config_get_string(&ptr, c.ptr, cname)
	if ret < 0 {
		return "", LastError()
	}

	return C.GoString(ptr), nil
}

func (c *Config) LookupBool(name string) (bool, error) {
	var out C.int
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.git_config_get_bool(&out, c.ptr, cname)
	if ret < 0 {
		return false, LastError()
	}

	return gobool(out), nil
}

func (c *Config) SetString(name, value string) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	ret := C.git_config_set_string(c.ptr, cname, cvalue)
	if ret < 0 {
		return LastError()
	}

	return nil
}

func (c *Config) SetInt32(name string, value int32) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.git_config_set_int32(c.ptr, cname, C.int32_t(value))
	if ret < 0 {
		return LastError()
	}

	return nil
}

func (c *Config) SetInt64(name string, value int64) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.git_config_set_int64(c.ptr, cname, C.int64_t(value))
	if ret < 0 {
		return LastError()
	}

	return nil
}

func (c *Config) SetBool(name string, value bool) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.git_config_set_bool(c.ptr, cname, cbool(value))
	if ret < 0 {
		return LastError()
	}

	return nil
}

func (c *Config) SetMultivar(name, regexp, value string) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cregexp := C.CString(regexp)
	defer C.free(unsafe.Pointer(cregexp))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	ret := C.git_config_set_multivar(c.ptr, cname, cregexp, cvalue)
	if ret < 0 {
		return LastError()
	}

	return nil
}

func (c *Config) Delete(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.git_config_delete_entry(c.ptr, cname)

	if ret < 0 {
		return LastError()
	}

	return nil
}

// OpenLevel creates a single-level focused config object from a multi-level one
func (c *Config) OpenLevel(parent *Config, level ConfigLevel) (*Config, error) {
	config := new(Config)
	ret := C.git_config_open_level(&config.ptr, parent.ptr, C.git_config_level_t(level))
	if ret < 0 {
		return nil, LastError()
	}

	return config, nil
}

// OpenOndisk creates a new config instance containing a single on-disk file
func OpenOndisk(parent *Config, path string) (*Config, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	config := new(Config)
	ret := C.git_config_open_ondisk(&config.ptr, cpath)
	if ret < 0 {
		return nil, LastError()
	}

	return config, nil
}

// Refresh refreshes the configuration to reflect any changes made externally e.g. on disk
func (c *Config) Refresh() error {
	ret := C.git_config_refresh(c.ptr)
	if ret < 0 {
		return LastError()
	}

	return nil
}

func (c *Config) Free() {
	runtime.SetFinalizer(c, nil)
	C.git_config_free(c.ptr)
}
