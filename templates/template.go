package templates

import (
	"os"
	"fmt"
	"path/filepath"
	"reflect"
	"text/template"
	"strings"
	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/utils"

	"golang.org/x/sys/unix"	
)

//Tmpl is a templating structure which implements the interface
type Tmpl struct{
	Path	string
}

//API is the templating interface for the configuration
type API interface{
	VersionExists(ver string)(error)
	CompileConfig(cfg *types.Config, dst string)(error)
}

//New returns a new Tmpl structure
func New(path string)(*Tmpl, error){
	if _, err := os.Stat(path); err!=nil{
		return nil,fmt.Errorf("Unable to create templating interface on path %s: %s", path, err)
	}
	return &Tmpl{Path: path}, nil
}

//VersionExists returns nil if a version exists, error otherwise
func (T *Tmpl) VersionExists(ver string)(error){
	_, err := os.Stat(filepath.Join(T.Path, ver))
	if os.IsNotExist(err){
		return fmt.Errorf("Template version %s does not exist in path %s: %s", ver, T.Path, err)
	}
	return err
}

//CompileConfig creates the configuration generated from a types.Config into a folder (e.g. in the runtime)
func (T *Tmpl) CompileConfig(cfg *types.Config, dst string)(error){
	//Check version exists
	if err := T.VersionExists(cfg.Version); err!=nil{
		return err
	}
	//check destination path is OK
	if unix.Access(dst, unix.F_OK) == nil{
		//check for read/write
		if e:=unix.Access(dst, unix.R_OK+unix.W_OK); e != nil{
			return fmt.Errorf("Cannot read/write in dst path %s: %s", dst, e)
		}
	}else{
		//create it
		if e:=os.Mkdir(dst, os.ModeDir|os.FileMode(0755)); e!=nil{
			return fmt.Errorf("Cannot create dst dir %s: %s", dst, e)
		}
	}
	//1. Does the conf parameters of Content exist in named files
	//1.1 Get the config json field name 
	field, ok := reflect.TypeOf(cfg.Content).Elem().FieldByName("Enb")
	if !ok {
		return fmt.Errorf("Config field 'Enb' not found")
	}
	fieldTag := field.Tag.Get("json") //this should return enb for Enb field
	//1.2 search the template path for "fieldTag".XXX.tmpl
	ext := fieldTag+".*"+"tmpl"
	match := utils.FindFileExt(ext, filepath.Join(T.Path, cfg.Version))
	if len(match) == 0{
		return fmt.Errorf("No template file for your config %s. The template should be in path %s and named %s.*.tmpl", fieldTag, filepath.Join(T.Path, cfg.Version), fieldTag)
	}
	//2. Now that everything is ok generate the damn thing
	t, err := template.ParseFiles(filepath.Join(T.Path, cfg.Version, match[0]))
	if err != nil{
		return fmt.Errorf("Unable to parse the template file: %s", err)
	}
	newConfFile := filepath.Join(dst, strings.TrimSuffix(match[0], ".tmpl"))
	f, err := os.OpenFile(newConfFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Unable to create new open new post: %s", err)
	}
	defer f.Close()
	//save the template in the conf dir of dst
	err = t.Execute(f, cfg.Content)
	if err != nil{
		return fmt.Errorf("Unable to execute template: %s", err)
	}

	return nil
}