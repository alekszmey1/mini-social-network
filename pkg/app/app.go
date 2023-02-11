package app

import (
	"log"

	"net/http"

	"homework.31/pkg/controller"
	"homework.31/pkg/repository"
	"homework.31/pkg/usecase"
)

func Run(port string) {
	repository := repository.NewRepository()
	usecase := usecase.NewUsecase(repository)
	controller := controller.NewController(usecase)
	mux := http.NewServeMux()

	mux.HandleFunc("/create", controller.CreateUser)
	mux.HandleFunc("/make_friends", controller.MakeFriends)
	mux.HandleFunc("/delete", controller.DeleteUser)
	mux.HandleFunc("/get_friends", controller.GetFriends)
	mux.HandleFunc("/put", controller.UpdateAge)

	go func() {
		log.Printf("запускаем сервер %s", port)
		http.ListenAndServe(port, mux)
	}()
}
