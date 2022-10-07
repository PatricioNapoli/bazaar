package utils

import "io/ioutil"

func ReadFile(path string) ([]byte, error) {
	dat, err := ioutil.ReadFile(path)
	return dat, err
}
