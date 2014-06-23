package motto

import (
    "os"
    "io/ioutil"
    "encoding/json"
)

func isDir(path string) (bool, error) {
    fi, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, nil
        }
        return false, err
    }

    return fi.IsDir(), nil
}

func isFile(path string) (bool, error) {
    fi, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, nil
        }
        return false, err
    }

    return !fi.IsDir(), nil
}

type packageInfo struct {
    Index string `json:"index"`
}

func parsePackageJsonIndex(path string) (string, error) {
    bytes, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }

    var info packageInfo
    err = json.Unmarshal(bytes, &info)
    if err != nil {
        return "", err
    }

    return info.Index, nil
}