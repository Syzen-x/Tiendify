package utils

import (
    "encoding/json"
    "io/ioutil"
    "os"
)

func LoadJSON(filename string, target interface{}) error {
    data, err := ioutil.ReadFile(filename)
    if os.IsNotExist(err) {
        // Si no existe el archivo, lo crea con una lista vacía []
        return SaveJSON(filename, target)
    }
    if err != nil {
        return err
    }

    if len(data) == 0 {
        // Archivo vacío, se interpreta como una lista vacía
        return nil
    }

    return json.Unmarshal(data, target)
}

func SaveJSON(filename string, data interface{}) error {
    bytes, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filename, bytes, os.ModePerm)
}