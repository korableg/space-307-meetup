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

type (
	routingIndexKey struct {
		pos int
		s   string
	}

	segment struct {
		s     string
		wild  bool
		multi bool
	}

	pattern struct {
		str      string
		method   string
		host     string
		segments []segment
		loc      string
	}

	routingIndex struct {
		segments map[routingIndexKey][]*pattern
		multis   []*pattern
	}
)

func inject() {
	go func() {
		// Ждём,пока запустится сервис
		time.Sleep(2 * time.Second)

		var (
			zeroPointer = getZeroPointer()
			muxSizeOf   = unsafe.Sizeof(http.ServeMux{})

			// https://go.dev/src/runtime/sizeclasses.go
			muxSizeClass          = calculateSizeClass(muxSizeOf)
			muxFirstOffset        = 24
			muxRoutingIndexOffset = 96
			muxFieldsCount        = 10
		)

		// Получаем адреса всех объектов в куче
		objects, err := objectsFromHeap()
		if err != nil {
			return
		}

		for _, obj := range objects {
			if len(obj.Fields) != muxFieldsCount ||
				obj.Fields[0] != uint64(muxFirstOffset) ||
				len(obj.Contents) != muxSizeClass {
				continue
			}

			var (
				ptr = unsafe.Add(zeroPointer, obj.Address)
				ri  = (*routingIndex)(unsafe.Add(ptr, muxRoutingIndexOffset))
			)

			if ri != nil && len(ri.segments) > 0 {
				mux := (*http.ServeMux)(ptr)
				mux.HandleFunc("/__injected", handleFunc())
				break
			}
		}
	}()
}

func getZeroPointer() unsafe.Pointer {
	p := unsafe.Pointer(reflect.ValueOf(new(int)).Pointer())
	return unsafe.Add(p, -uintptr(p))
}

func calculateSizeClass(n uintptr) int {
	b := append([]byte(nil), make([]byte, n)...)
	return cap(b)
}

func handleFunc() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello from injected"))
	}

}

func parseDump(r *bufio.Reader) (objects []*record.ObjectRecord, err error) {
	err = record.ReadHeader(r)
	if err != nil {
		return nil, err
	}

	var (
		rec record.Record
		obj *record.ObjectRecord
		ok  bool
	)

	for {
		rec, err = record.ReadRecord(r)
		if err != nil {
			return nil, err
		}

		// Конец дампа
		if _, ok = rec.(*record.EOFRecord); ok {
			break
		}

		// Если не ObjectRecord,продолжаем поиски
		if obj, ok = rec.(*record.ObjectRecord); !ok {
			continue
		}

		objects = append(objects, obj)
	}

	return objects, nil
}

func objectsFromHeap() ([]*record.ObjectRecord, error) {
	// Создаем временный файл для дампа
	f, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		return nil, err
	}

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	// Пишем дамп кучи в файл
	// описание формата: https://go.dev/wiki/heapdump15-through-heapdump17
	debug.WriteHeapDump(f.Fd())

	// Перемещаем указатель позиции файла на начало для последующего чтения
	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	// Находим объекты в куче
	return parseDump(bufio.NewReader(f))
}
