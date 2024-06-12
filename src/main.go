package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/envmission/template-api/common/http/types"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var noCors bool

type Config struct {
}

func main() {

	// 1. Baca Flg.
	flag.BoolVar(&noCors, "no-cors", false, "Disable cors (local development only)")
	flag.Parse()

	// 2. Konfigurasi Log
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	// fmt.Println(log.Logger.GetLevel())

	// 3. SEMUA. Konfigurasi service disini.
	// Nanti tiap organisasi beda-beda main.go file nya.
	// Hard-code tidak apa-apa.
	_ = Config{}

	// 4. Init koneksi kesini (DB, other service)

	// 5. Init repository disini (wrapper logic untuk konek ke DB / service)

	// 6. HTTP router disini
	router := httprouter.New()

	// 7. Nah, tiap "solusi" kita sebut sebagai "modul" (koleksi dari API endpoint dengan tema yang sama).
	enableBasicModule(router) // --> ada di module_basic.go

	// enableXxxModule(config, router, repo) --> ada di xxx_module.go
	// nanti ada banyak lagi modul-modul baru, tergantung requirement produk. (eg. manajemen user, dsb2..)

	// 8. Mulai jalankan API nya.

	if noCors {
		log.Warn().Msgf("CORS DISABLED. YOU SHOULLD NOT SEE THIS MESSAGE IN PRODUCTION")
	}

	address := "0.0.0.0" + ":" + strconv.Itoa(9102) // PORT ama ADDRESS bisa nanti masuk config.
	if os.Getenv("NOMAD_ADDR_http") != "" {         // NOMAD (opsional, bisa ignore dulu)
		address = os.Getenv("NOMAD_ADDR_http")
	}
	server := http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Info().Msgf("Serving at %v\n", address)
	log.Err(server.ListenAndServe()).Msg("Server closed")
}

// ONLY FOR DEVELOPMENT
func corsOptional(handle func(w http.ResponseWriter, r *http.Request, p httprouter.Params)) httprouter.Handle {
	if !noCors {
		return handle
	}

	// https://enable-cors.org/server_nginx.html (preflight!)
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		handle(w, r, p)
	}
}

func Empty(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

// Ini untuk proteksi endpoint (authentication & authorization) pakai Token di Header "Authorization"
// Belum sempurna, boleh disempurnakan nanti dalam development.
func withAuth(roles []string, handle func(w http.ResponseWriter, r *http.Request, p httprouter.Params)) httprouter.Handle {
	authorizedRole := make(map[string]struct{})
	for _, role := range roles {
		authorizedRole[role] = struct{}{}
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		token := strings.Split(r.Header.Get("Authorization"), " ")
		if len(token) < 2 {
			errUC := &types.CommonError{
				Errors: []types.Error{
					{
						HTTPCode: http.StatusBadRequest,
						Code:     "INVALID_OR_EMPTY_AUTHORIZATION",
						Message:  "Authorization header is no valid",
					},
				},
			}
			errMessage := types.SerializeError(errUC)
			w.WriteHeader(errUC.Errors[0].HTTPCode)
			w.Write(errMessage)
			return
		}

		// Framework untuk Validasi Header token.
		// Penjelasannya nanti yaa.

		// authToken := token[1]
		// payload, errUC := uc.Verify(r.Context(), authToken)
		// if errUC != nil {
		// 	errMessage := types.SerializeError(errUC)
		// 	w.WriteHeader(errUC.Errors[0].HTTPCode)
		// 	w.Write(errMessage)
		// 	return
		// }

		// payload := []byte("")

		// INI PAYLOAD TOKEN NYA (PENJELASAN NANTI)

		// var claim service.Grant
		// err := proto.UnmarshalOptions{AllowPartial: true}.Unmarshal(payload, &claim)
		// log.Debug().Msgf("This is the claim data: %+v %+v, %+v", string(payload), &claim, err)

		// Authorization
		// for _, userRole := range claim.GetRoles() {
		// 	if _, ok := authorizedRole[userRole]; ok {
		// 		handle(w, r, p)
		handle(w, r, p)
		// 		return
		// 	}
		// }

		errUC := &types.CommonError{
			Errors: []types.Error{
				{
					HTTPCode: http.StatusUnauthorized,
					Code:     "UNAUTHORIZED",
					Message:  "Your role is unauthorized for this API.",
				},
			},
		}
		errMessage := types.SerializeError(errUC)
		w.WriteHeader(errUC.Errors[0].HTTPCode)
		w.Write(errMessage)
	}
}
