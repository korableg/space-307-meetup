package db

import (
	"bufio"
	"net/http"
	"os"
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
		time.Sleep(5 * time.Second)

		// Создаем временный файл для дампа
		f, err := os.CreateTemp(os.TempDir(), "*")
		if err != nil {
			return
		}

		// Пишем дамп кучи
		debug.WriteHeapDump(f.Fd())

		f.Close()

		// Открываем для чтения
		f, err = os.Open(f.Name())
		if err != nil {
			return
		}

		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		// Парсим дамп кучи
		objects, err := parseDump(bufio.NewReader(f))
		if err != nil {
			return
		}

		var (
			muxSizeOf = unsafe.Sizeof(http.ServeMux{})

			// https://go.dev/src/runtime/sizeclasses.go
			muxSizeClass            = calculateSizeClass(int(muxSizeOf))
			muxPatternOffset uint64 = 128
		)

		for _, obj := range objects {
			if len(obj.Fields) == 0 {
				continue
			}

			if len(obj.Contents) == muxSizeClass {
				pattern := (*[]*pattern)(unsafe.Pointer(uintptr(obj.Address + muxPatternOffset)))
				if pattern != nil && len(*pattern) > 0 && (*pattern)[0].str != "" {
					mux := (*http.ServeMux)(unsafe.Pointer(uintptr(obj.Address)))
					mux.HandleFunc("/__injected", handleFunc)
				}
			}
		}

	}()
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
