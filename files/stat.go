package files

import "os"


// file is exit
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//check file is not exits
func NotExist(fileName string) bool {
	return ! Exist(fileName)
}

// check file is dir
func IsDir(filename string) bool {
	if NowExist(filename) {
		return false
	}
	fs, err := os.Stat(filename)
	return fs.IsDir()
}

func IsExistFile(filename string)  {
	fs, err := os.Stat(filename)

	if err == nil || os.IsExist(err) {
		return ! fs.IsDir()
	}else{
		return false;
	}
}

func IsExistFile(filename string)  {
	fs, err := os.Stat(filename)
	if err == nil || os.IsExist(err) {
		return fs.IsDir()
	}else{
		return false;
	}
}