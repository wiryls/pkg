package config

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/wiryls/pkg/errors/cerrors"
	"github.com/wiryls/pkg/errors/wrap"
)

// Save writes the configuration to target path.
func (s Serializers) Save(src interface{}, dst string) (err error) {
	err = cerrors.TestNilArgumentIfNoErr(err, src, "src interface{}")
	if err == nil {
		err = s.toFile(dst, src)
	}
	return
}

// Load a configuration form target path.
func (s Serializers) Load(dst interface{}, src string) (err error) {
	err = cerrors.TestNilArgumentIfNoErr(err, dst, "dst interface{}")
	if err == nil {
		err = s.fromFile(src, dst)
	}
	return
}

// LoadSome reads a list of files and finds the first valid one as configuration.
func (s Serializers) LoadSome(dst interface{}, list []string) (path string, err error) {

	err = cerrors.TestNilArgumentIfNoErr(err, list, "list []string")
	err = cerrors.TestNilArgumentIfNoErr(err, dst, "dst interface{}")

	if err == nil {
		msg := strings.Builder{}
		for _, p := range list {
			if err := s.fromFile(p, dst); err == nil {
				path = p
				break
			} else {
				if msg.Len() > 0 {
					msg.WriteString("; ")
				}
				msg.WriteString(err.Error())
			}
		}

		switch {
		case path != "":
			err = nil
		case msg.Len() == 0:
			err = wrap.Message(ErrReading, "no files")
		default:
			err = wrap.Message(ErrReading, msg.String())
		}
	}

	return
}

// LoadOrSaveSome trys to read a config from some paths.
// If failed, save dst to the first valid path.
func (s Serializers) LoadOrSaveSome(data interface{}, list []string) (path string, find bool, err error) {

	if err == nil {
		path, err = s.LoadSome(data, list)
		find = err == nil
	}

	if err != nil {
		for _, path := range list {
			if err := s.Save(data, path); err == nil {
				return path, false, nil
			}
		}
	}

	return
}

func (s Serializers) fromFile(path string, target interface{}) (err error) {
	var buf io.ReadCloser

	if err == nil {
		_, err = os.Stat(path)
		err = wrap.MessageAliasStack(err, "cannot find `"+path+"`", ErrReading, 0)
	}

	if err == nil {
		buf, err = os.Open(path)
		err = wrap.MessageAliasStack(err, "cannot open `"+path+"`", ErrReading, 0)
	}

	if err == nil {
		defer buf.Close()
		err = s.fromReader(path, buf, target)
	}

	return
}

func (s Serializers) toFile(path string, src interface{}) (err error) {
	var buf io.WriteCloser

	if err == nil {
		buf, err = os.Create(path)
		err = wrap.MessageAliasStack(err, "cannot create `"+path+"`", ErrWriting, 0)
	}

	if err == nil {
		defer buf.Close()
		err = s.toWriter(path, buf, src)
	}

	return
}

func (s Serializers) fromReader(name string, r io.Reader, dst interface{}) (err error) {

	err = cerrors.TestNilArgumentIfNoErr(err, r, "r io.Reader")
	err = cerrors.TestNilArgumentIfNoErr(err, dst, "dst interface{}")

	var bin []byte
	if err == nil {
		bin, err = ioutil.ReadAll(r)
		err = wrap.MessageAliasStack(err, "cannot read `"+name+"`", ErrReading, 0)
	}

	return s.fromBytes(name, bin, dst)
}

func (s Serializers) toWriter(name string, w io.Writer, src interface{}) (err error) {

	err = cerrors.TestNilArgumentIfNoErr(err, w, "w io.Writer")
	err = cerrors.TestNilArgumentIfNoErr(err, src, "src interface{}")

	var bin []byte
	if err == nil {
		bin, err = s.toBytes(name, src)
	}

	if err == nil {
		_, err = w.Write(bin)
		err = wrap.MessageAliasStack(err, "cannot write `"+name+"`", ErrWriting, 0)
	}

	return
}

func (s Serializers) fromBytes(name string, bin []byte, dst interface{}) (err error) {

	err = cerrors.TestNilArgumentIfNoErr(err, bin, "bin []byte")
	err = cerrors.TestNilArgumentIfNoErr(err, dst, "dst interface{}")

	if err == nil {
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
		msg := "failed to encode " + ext

		for _, s := range s {
			if s.CanUnmarshal(ext, bin) {
				err = s.Unmarshal(bin, dst)
				err = wrap.MessageAliasStack(err, msg, ErrDecoding, 0)
				return
			}
		}
	}

	return wrap.MessageStack(ErrDecoding, "file not support `"+name+"`", 0)
}

func (s Serializers) toBytes(name string, src interface{}) (bin []byte, err error) {

	err = cerrors.TestNilArgumentIfNoErr(err, src, "src interface{}")

	if err == nil {
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
		msg := "failed to encode " + ext

		for _, s := range s {
			if s.CanMarshal(ext, src) {
				bin, err = s.Marshal(src)
				err = wrap.MessageAliasStack(err, msg, ErrEncoding, 0)
				return
			}
		}
	}

	return nil, wrap.MessageStack(ErrEncoding, "file not support `"+name+"`", 0)
}
