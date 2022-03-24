package interfaces

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"news-dd/application"
	"news-dd/domain"
	"strconv"
	"strings"
	"time"
)

type authConfig struct {
	Usr string
	Pwd string
}

func Run(port int) error {
	log.Printf("Server running at http://localhost:%d/", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), Routes())
}

//Basic auth check
func BasicAuth(h httprouter.Handle, user, pass []byte) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		const basicAuthPrefix string = "Basic "

		// Get the Basic Authentication credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 && bytes.Equal(pair[0], user) && bytes.Equal(pair[1], pass) {
					// Delegate request to the given handle
					h(w, r, ps)
					return
				}
			}
		}

		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

// Routes returns the initialized router
func Routes() *httprouter.Router {
	var auth *authConfig
	if _, err := toml.DecodeFile("config.toml", &auth); err != nil {
		log.Fatal(err)
	}

	user := []byte(auth.Usr)
	pass := []byte(auth.Pwd)

	r := httprouter.New()

	// Index Route
	r.GET("/", index)
	r.GET("/api/v1", index)

	// Topic Route
	r.GET("/api/v1/topic", BasicAuth(getAllTopic, user, pass))
	r.GET("/api/v1/topic/:topic_id", BasicAuth(getTopic, user, pass))
	r.POST("/api/v1/topic", BasicAuth(createTopic, user, pass))
	r.DELETE("/api/v1/topic/:topic_id", BasicAuth(removeTopic, user, pass))
	r.PUT("/api/v1/topic/:topic_id", BasicAuth(updateTopic, user, pass))

	// Tags Route
	r.GET("/api/v1/tags", BasicAuth(getAllTags, user, pass))
	r.GET("/api/v1/tags/:tags_id", BasicAuth(getTags, user, pass))
	r.POST("/api/v1/tags", BasicAuth(createTags, user, pass))
	r.DELETE("/api/v1/tags/:tags_id", BasicAuth(removeTags, user, pass))
	r.PUT("/api/v1/tags/:tags_id", BasicAuth(updateTags, user, pass))

	// News Route
	r.POST("/api/v1/news/:topic_id/:status", BasicAuth(getAllNews, user, pass))
	r.GET("/api/v1/news/:news_id", BasicAuth(getNews, user, pass))
	r.POST("/api/v1/news", BasicAuth(createNews, user, pass))
	r.DELETE("/api/v1/news/:news_id", BasicAuth(removeNews, user, pass))
	r.PUT("/api/v1/news/:news_id", BasicAuth(updateNews, user, pass))

	return r
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	JSON(w, http.StatusOK, "GO DD BAREKSA TEST")
}

// =============================
//    TOPIC
// =============================

func getTopic(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	topicID, err := strconv.Atoi(ps.ByName("topic_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	topic, err := application.GetTopic(topicID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, topic)
}

func getAllTopic(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	topics, err := application.GetAllTopic()
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, topics)
}

func createTopic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	type payload struct {
		TopicName string `json:"topic"`
	}
	var p payload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	err = application.AddTopic(p.TopicName)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	JSON(w, http.StatusCreated, nil)
}

func removeTopic(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	topicID, err := strconv.Atoi(ps.ByName("topic_id"))
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = application.RemoveTopic(topicID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}

func updateTopic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var p domain.Topic
	err := decoder.Decode(&p)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
	}

	topicID, err := strconv.Atoi(ps.ByName("topic_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	err = application.UpdateTopic(p, topicID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}

// =============================
//    TAGS
// =============================

func getTags(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	tagsID, err := strconv.Atoi(ps.ByName("tags_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	tags, err := application.GetTags(tagsID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, tags)
}

func getAllTags(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	tags, err := application.GetAllTags()
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, tags)
}

func createTags(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type payload struct {
		TagName string `json:"tag"`
	}
	var p payload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	err = application.AddTags(p.TagName)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	JSON(w, http.StatusCreated, nil)
}

func removeTags(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	tagsID, err := strconv.Atoi(ps.ByName("tags_id"))
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = application.RemoveTags(tagsID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}

func updateTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var p domain.Tags
	err := decoder.Decode(&p)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
	}

	tagsID, err := strconv.Atoi(ps.ByName("tags_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	err = application.UpdateTags(p, tagsID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}

// =============================
//    NEWS
// =============================

func getNews(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	newsId, err := strconv.Atoi(ps.ByName("news_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	news, err := application.GetNews(newsId)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	topic, err := application.GetTopic(news.Topic)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	newsRespond := domain.NewsRespond{
		Id:         news.Id,
		Title:      news.Title,
		Content:    news.Content,
		Status:     news.Status,
		Topic:      topic.Topic,
		Tags:       news.Tags,
		CreateTime: news.CreateTime,
		UpdateTime: news.UpdateTime,
	}

	JSON(w, http.StatusOK, newsRespond)
}

func getAllNews(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	topicId, err := strconv.Atoi(ps.ByName("topic_id"))
	if err != nil {
		topicId = 0
	}
	status := ps.ByName("status")
	if status != "draft" && status != "deleted" && status != "publish" {
		status = ""
	}

	news, err := application.GetAllNews(topicId, status)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	newsRespond := []domain.NewsRespond{}

	for _, v := range news {
		resp := domain.NewsRespond{
			Id:         v.Id,
			Title:      v.Title,
			Content:    v.Content,
			Status:     v.Status,
			Tags:       v.Tags,
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
		}

		topic, _ := application.GetTopic(v.Topic)
		resp.Topic = topic.Topic

		newsRespond = append(newsRespond, resp)
	}

	JSON(w, http.StatusOK, newsRespond)
}

func createNews(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p domain.News

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}
	p.CreateTime = time.Now()

	err = application.AddNews(p)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	JSON(w, http.StatusCreated, nil)
}

func removeNews(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	newsID, err := strconv.Atoi(ps.ByName("news_id"))
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = application.RemoveNews(newsID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}

func updateNews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var p domain.News
	err := decoder.Decode(&p)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
	}

	newsID, err := strconv.Atoi(ps.ByName("news_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	err = application.UpdateNews(p, newsID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}
