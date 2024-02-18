package customerror

import (
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					_, err := w.Write(ErrNotFound.Marshal())
					if err != nil {
						return
					}
					return
				}

				err = err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write(ErrNotFound.Marshal())
				if err != nil {
					return
				}
				return
			}

			w.WriteHeader(http.StatusTeapot)
			_, err := w.Write(systemError(err).Marshal())
			if err != nil {
				return
			}
		}
	}
}
