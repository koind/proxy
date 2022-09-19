package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koind/proxy/internal/domain/repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTP сервер
type HTTPServer struct {
	http.Server
	router *gin.Engine
	domain string
	repo   repository.RequestRepositoryInterface
	client *http.Client
}

// Возвращает новый HTTP сервер
func NewHTTPServer(repo repository.RequestRepositoryInterface, domain string) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.New()
	c := &http.Client{Timeout: 5 * time.Second}
	hs := HTTPServer{router: r, domain: domain, repo: repo, client: c}

	hs.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	hs.router.POST("/", hs.sendHandle)
	hs.router.GET("/ping", hs.pingHandle)
	hs.router.GET("/history", hs.historyHandle)

	http.Handle("/", r)

	return &hs
}

// Запускает HTTP сервер
func (s *HTTPServer) Start() error {
	return http.ListenAndServe(s.domain, s.router)
}

// Отправляет запрос и сохраняет данные
// @Summary Отправляет запрос и сохраняет данные
// @Description Отправляет запрос и сохраняет данные
// @Tags api
// @ID send-handle
// @Accept  json
// @Produce  json
// @Param data body repository.Request true "входные данные"
// @Success 200 {object} repository.Response
// @Failure 400,404 {object} handler.response
// @Failure 500 {object} handler.response
// @Router / [post]
func (s *HTTPServer) sendHandle(c *gin.Context) {
	inputReq := repository.Request{}

	if err := json.NewDecoder(c.Request.Body).Decode(&inputReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError("error decoding the request"))
		return
	}

	outputReq, err := http.NewRequestWithContext(c, inputReq.Method, inputReq.URL, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError("error building the request"))
		return
	}

	setRequestHeaders(outputReq, inputReq.Headers)

	resp, err := s.client.Do(outputReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError("error doing the request"))
		return
	}
	defer resp.Body.Close()

	rsp := repository.Response{
		ID:      uuid.New().String(),
		Status:  resp.Status,
		Headers: headers(resp.Header),
		Length:  resp.ContentLength,
	}
	row := repository.Row{ID: rsp.ID, Request: inputReq, Response: rsp}

	if _, err := s.repo.Create(c, row); err != nil {
		c.JSON(http.StatusBadRequest, responseError(err.Error()))
	}

	c.JSON(http.StatusOK, rsp)
}

// Ping запрос для проверки сервиса
// @Summary Ping запрос для проверки сервиса
// @Description Ping запрос для проверки сервиса
// @Tags api
// @ID pint-handle
// @Success 200
// @Router /ping [get]
func (s *HTTPServer) pingHandle(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

// Возвращет все записи
// @Summary Возвращет все записи
// @Description Возвращет все записи
// @Tags api
// @ID get-all-history-handle
// @Produce  json
// @Success 200 {array} repository.Response
// @Failure 400,404 {object} handler.response
// @Failure 500 {object} handler.response
// @Router /history [get]
func (s *HTTPServer) historyHandle(c *gin.Context) {
	list, err := s.repo.GetAll(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError(err.Error()))
	}

	c.JSON(http.StatusOK, list)
}

func setRequestHeaders(destReq *http.Request, srcHeaders map[string]string) {
	for k, v := range srcHeaders {
		destReq.Header.Set(k, v)
	}
}

func headers(h http.Header) map[string]string {
	hs := make(map[string]string)
	for k := range h {
		hs[k] = h.Get(k)
	}
	return hs
}

type response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func responseError(message string) response {
	return response{Status: "error", Message: message}
}
