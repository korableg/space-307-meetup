package db

import (
	"bufio"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"time"
	"unsafe"

	"github.com/burntcarrot/heaputil/record"
)

type pattern struct {
	str    string
	method string
	host   string
}

func inject() {
	go func() {
		// Ждём,пока запустится сервис
		time.Sleep(2 * time.Second)

		// Создаем временный файл для дампа
		f, err := os.CreateTemp(os.TempDir(), "*")
		if err != nil {
			return
		}

		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		// Пишем дамп кучи
		debug.WriteHeapDump(f.Fd())

		// Перемещаем указатель позиции файла на начало
		_, err = f.Seek(0, 0)
		if err != nil {
			return
		}

		// Парсим дамп
		objects, err := parseDump(bufio.NewReader(f))
		if err != nil {
			return
		}

		var (
			zeroPointer = getZeroPointer()
			muxSizeOf   = unsafe.Sizeof(http.ServeMux{})

			// https://go.dev/src/runtime/sizeclasses.go
			muxSizeClass     = calculateSizeClass(int(muxSizeOf))
			muxPatternOffset = 128
		)

		for _, obj := range objects {
			if len(obj.Fields) == 0 {
				continue
			}

			if len(obj.Contents) == muxSizeClass {
				ptr := unsafe.Add(zeroPointer, obj.Address)
				pattern := (*[]*pattern)(unsafe.Add(ptr, muxPatternOffset))
				if pattern != nil && len(*pattern) > 0 && (*pattern)[0].str != "" {
					mux := (*http.ServeMux)(ptr)
					mux.HandleFunc("/__injected", handleFunc)
				}
			}
		}

	}()
}

func getZeroPointer() unsafe.Pointer {
	p := unsafe.Pointer(reflect.ValueOf(new(int)).Pointer())
	return unsafe.Add(p, -uintptr(p))
}

func calculateSizeClass(n int) int {
	b := append([]byte(nil), make([]byte, n)...)
	return cap(b)
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello from injected!"))
}

func parseDump(r *bufio.Reader) ([]*record.ObjectRecord, error) {
	err := record.ReadHeader(r)
	if err != nil {
		return nil, err
	}

	var objects []*record.ObjectRecord

	for {
		r, err := record.ReadRecord(r)
		if err != nil {
			return nil, err
		}

		_, isEOF := r.(*record.EOFRecord)
		if isEOF {
			break
		}

		obj, ok := r.(*record.ObjectRecord)
		if !ok {
			continue
		}

		objects = append(objects, obj)
	}

	return objects, nil
}
