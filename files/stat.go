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
	if NotExist(filename) {
		return false
	}
	fs, _ := os.Stat(filename)
	return fs.IsDir()
}

func IsExistFile(filename string) bool  {
	fs, err := os.Stat(filename)

	if err == nil || os.IsExist(err) {
		return ! fs.IsDir()
	}else{
		return false;
	}
}

func IsExistDir(filename string) bool {
	fs, err := os.Stat(filename)
	if err == nil || os.IsExist(err) {
		return fs.IsDir()
	}else{
		return false;
	}
}
