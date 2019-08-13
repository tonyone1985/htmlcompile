package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
)
func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
func getPages(dir string,reladir string)([]string, []string){
	files, _ := ioutil.ReadDir(dir)
	fss := make([]string,0)
	odir :=make([]string,0)
	for _,f :=range files{
		fp := path.Join(dir,f.Name())
		if f.IsDir(){
			cp,cdir:= getPages(fp,path.Join(reladir,f.Name()))
			fss=append(fss, cp...)
			odir=append(odir,cdir...)

			continue
		}
		fss=append(fss, fp)
		odir=append(odir,reladir)
	}
	return fss,odir
}
func main()  {
	//t,e:=template.ParseGlob("F:/src/go/webbase/static/pages/*")
	root:="./"

	if len(os.Args)>1{
		root = os.Args[1]
	}
	release := path.Join(root,"release")

	if len(os.Args)>2{
		release = os.Args[2]
	}
	otherdirs, _ := ioutil.ReadDir(root)
	for _,d :=range otherdirs{
		if d.Name() == "tmpl" ||d.Name() == "pages" || d.Name() == "release"{
			continue
		}
		if !d.IsDir(){
			continue
		}
		files,dirs := getPages(path.Join(root,d.Name()),path.Join(release,d.Name()))
		for i:=0;i<len(files);i++{
			_,fname:= path.Split(files[i])
			CopyFile(files[i],path.Join(dirs[i],fname))
		}
	}
	t,_ := getPages(path.Join(root,"tmpl"),release)
	files,dirs := getPages(path.Join(root,"pages"),release)
	t=append(files,t...)
	tm,_:= template.ParseFiles(t...)
	for i:=0;i<len(files);i++{
		if !IsExist( dirs[i]){
			os.MkdirAll(dirs[i],os.ModePerm)
		}
		_,fname:= path.Split(files[i])
		finame := path.Join(dirs[i],fname)
		f,_ := os.Create(finame)
		tm.ExecuteTemplate(f,fname,nil)
		f.Close()
	}


	//t.ParseFiles()

	//t,e:=template.ParseGlob("*.html")
	//
	//
	////t2,_:=template.ParseGlob("F:/src/go/webbase/static/pages/*.tmpl")
	////t,e:= template.ParseFiles("F:/src/go/webbase/static/pages/*.html","F:/src/go/webbase/static/pages/*.tmpl")
	////t,e:= template.ParseFiles("F:/src/go/webbase/static/pages/mydetails.html")
	////f,_ := os.Create("t.txt")
	////t.Execute(f,map[string]string{})
	////f.Close()
	//if e!=nil{
	//	fmt.Println("[err]",e)
	//}
	//reg := regexp.MustCompile(`[\d\D]+\.tmpl`)
	//for _,v:=range t.Templates(){
	//	//fmt.Println(v.Name())
	//	if v.Name() == "mydetails.tmpl"{
	//		fmt.Println("aaa")
	//	}
	//
	//	if reg.MatchString(v.Name() ){
	//	//if v.Name() == "mydetails.html" || v.Name()=="login.html"{
	//	//	t.ExecuteTemplate(os.Stdout,v.Name(),nil)
	//
	//	}
	//}


}




func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()

	dest_dir,_ := path.Split(dst)
	b := IsExist(dest_dir)
	if b == false {
		err := os.MkdirAll(dest_dir, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			fmt.Println(err)
		}
	}
	dstFile, err := os.Create(dst)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}
