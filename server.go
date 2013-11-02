package main

// #include <stddef.h>
// extern char _binary_assets_zip_start[];
// extern char _binary_assets_zip_end[];
// int resource_size() {return _binary_assets_zip_end - _binary_assets_zip_start;}
// char* resource() {return _binary_assets_zip_start;}
// #cgo LDFLAGS: -L . -lassets
import "C"
import "unsafe"

import (
    "net/http"
    "fmt"
    "bytes"
    "io/ioutil"
    "archive/zip"
    "log"
    "github.com/gorilla/mux"
)

type cache map[string] []byte 

var fileCache cache

// Pick up assets from the C library
func extractAssets() []byte {
    size := C.resource_size();
    fmt.Println(size)
    bytes := C.GoBytes(unsafe.Pointer(C.resource()), C.resource_size())
    return bytes
}

func readArchive(rawArchive []byte) cache  {
    fcache := make(map[string] []byte)

    // r, err := zip.OpenReader(archive)
    r, err := makeZipReader(rawArchive)
    if err != nil {
            log.Fatal(err)
    }
    // defer r.Close()

    for _, f := range r.File {
        log.Printf("Found file: %s", f.Name)
        rc, err := f.Open()
        if err != nil {
            log.Println("cannot open")
            log.Fatal(err)
        }
        defer rc.Close()

        // _, err = io.Copy(os.Stdout, rc)
        content, err := ioutil.ReadAll(rc)
        if err != nil {
            log.Println("cannot read file")
            log.Fatal(err)
        }
        fcache[f.Name] = content
    }
    return fcache
}

func makeZipReader(buffer []byte) (*zip.Reader, error) {
    reader := bytes.NewReader(buffer)
    r, err := zip.NewReader(reader, int64(len(buffer))) 
    if err != nil {
        return nil, err
    }
    return r, nil
}

func assetHandler(writer http.ResponseWriter, request *http.Request) {
    path :=request.URL.Path[1:]
    // log.Printf("Path: %s\n", path)
    content, found := fileCache[path]
    if found {
        writer.Write(content)
    } else {
        fmt.Fprintf(writer, "Not found")
    }
}

func main() {
    rawAssets := extractAssets();

    //fileCache = readArchive("assets.zip")
    fileCache = readArchive(rawAssets);

    r := mux.NewRouter()
    r.HandleFunc("/assets/{path:.*}", assetHandler)
    http.ListenAndServe(":3000", r)
}
