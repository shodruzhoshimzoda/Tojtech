package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/product"
)







type ProductHandler struct {
	usc 	*usecase.ProductUsecase	
	log 	*slog.Logger
}


func NewProductHandler(usc *usecase.ProductUsecase, log *slog.Logger)  *ProductHandler{
	return &ProductHandler{
		usc: usc,
		log: log,
	}
}

func (h *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {

	const op = "ProductHandler.GetProductHandler"

	id := chi.URLParam(r, "id")
	
	

	idint, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		h.log.Error("ERROR parsing id: ", slog.String("err", err.Error()), slog.String("op",op))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w,r,map[string]string{"error":"invalid id"})
		return
	}

	product, err := h.usc.GetProduct(r.Context(), idint)
	if err != nil {
		
		if errors.Is(err, usecase.ErrProductNotFound) {
			h.log.Error("Error getting product", slog.String("err", err.Error()), slog.String("op",op))
			render.Status(r, http.StatusNotFound)
			render.JSON(w,r,map[string]string{"error":"product not found"})
			return
		}

		
		h.log.Error("failed to get product", slog.String("err", err.Error()), slog.String("op", op))
        render.Status(r, http.StatusInternalServerError)
        render.JSON(w, r, map[string]string{"error": "internal server error"})
        return
		
	}

	h.log.Info("get prouduct: ", slog.String("id", id))
	render.JSON(w, r, product)
	
}